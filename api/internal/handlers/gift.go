package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/centrifugo"
	"github.com/huya_live/api/pkg/response"
)

type GiftHandler struct{}

func NewGiftHandler() *GiftHandler {
	return &GiftHandler{}
}

type SendGiftRequest struct {
	RoomID    string `json:"room_id" binding:"required"`
	GiftID    int    `json:"gift_id" binding:"required"`
	GiftCount int    `json:"gift_count" binding:"required,min=1,max=99"`
}

func (h *GiftHandler) SendGift(c *gin.Context) {
	userID := c.GetString("user_id")
	var req SendGiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if req.GiftCount < 1 || req.GiftCount > 99 {
		response.BadRequest(c, "gift_count must be between 1 and 99")
		return
	}

	var room models.LiveRoom
	var relay models.RelayStream
	roomFound := false
	isRelay := false

	if err := repository.DB.Where("id = ? AND status = ?", req.RoomID, "live").First(&room).Error; err == nil {
		roomFound = true
	} else {
		if err := repository.DB.Where("id = ? AND status = ?", req.RoomID, "running").First(&relay).Error; err == nil {
			roomFound = true
			isRelay = true
		}
	}

	if !roomFound {
		response.BadRequest(c, "room not found or not live")
		return
	}

	if !isRelay && room.StreamerID.String() == userID {
		response.BadRequest(c, "cannot send gift to yourself")
		return
	}

	if isRelay && relay.ID.String() == userID {
		response.BadRequest(c, "cannot send gift to yourself")
		return
	}

	var gift models.Gift
	if err := repository.DB.Where("id = ? AND is_active = ?", req.GiftID, true).First(&gift).Error; err != nil {
		response.BadRequest(c, "gift not found")
		return
	}

	if gift.MinLevelRequired > 1 {
		var user models.User
		if err := repository.DB.First(&user, "id = ?", userID).Error; err == nil {
			if user.Level < gift.MinLevelRequired {
				response.BadRequest(c, "level not high enough to send this gift")
				return
			}
		}
	}

	var user models.User
	if err := repository.DB.First(&user, "id = ?", userID).Error; err != nil {
		response.BadRequest(c, "user not found")
		return
	}

	totalCost := gift.CoinPrice * req.GiftCount
	if user.CoinBalance < totalCost {
		response.BadRequest(c, "insufficient coins")
		return
	}

	var streamer models.Streamer
	var streamerID uuid.UUID
	if isRelay {
		streamerID = relay.ID
	} else {
		streamerID = room.StreamerID
		if err := repository.DB.Where("user_id = ?", room.StreamerID).First(&streamer).Error; err != nil {
			response.BadRequest(c, "streamer not found")
			return
		}
	}

	tx := repository.DB.Begin()
	if tx.Error != nil {
		response.Fail(c, "failed to start transaction")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user.CoinBalance -= totalCost
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		response.Fail(c, "failed to deduct coins")
		return
	}

	coinTx := models.CoinTransaction{
		UserID:       user.ID,
		Amount:       -totalCost,
		BalanceAfter: user.CoinBalance,
		Type:         "gift",
		Description:  "Send gift to room " + room.ID.String(),
	}
	if err := tx.Create(&coinTx).Error; err != nil {
		tx.Rollback()
		response.Fail(c, "failed to record transaction")
		return
	}

	levelConfig := models.LevelConfig{}
	repository.DB.First(&levelConfig, "level = ?", user.Level)
	bonusMultiplier := levelConfig.BonusMultiplier
	if bonusMultiplier == 0 {
		bonusMultiplier = 1.0
	}
	loyaltyPoints := int64(float64(totalCost) * bonusMultiplier)

	if isRelay {
		response.Success(c, gin.H{
			"message":           "gift sent to relay stream (no streamer revenue)",
			"gift_name":         gift.Name,
			"gift_count":        req.GiftCount,
			"total_cost":        totalCost,
			"remaining_balance": user.CoinBalance,
		})
		return
	}

	giftTx := models.GiftTransaction{
		SenderID:            user.ID,
		ReceiverID:          streamerID,
		RoomID:              room.ID,
		GiftID:              gift.ID,
		GiftCount:           req.GiftCount,
		CoinAmount:          totalCost,
		LoyaltyPointsGained: loyaltyPoints,
		UserLevelAtSend:     user.Level,
		BonusMultiplier:     bonusMultiplier,
	}
	if err := tx.Create(&giftTx).Error; err != nil {
		tx.Rollback()
		response.Fail(c, "failed to record gift")
		return
	}

	if !isRelay {
		streamer.TotalRevenue += int64(totalCost)
		if err := tx.Save(&streamer).Error; err != nil {
			tx.Rollback()
			response.Fail(c, "failed to update streamer revenue")
			return
		}

		var fanRelation models.FanRelation
		if err := tx.Where("user_id = ? AND streamer_id = ?", userID, room.StreamerID).First(&fanRelation).Error; err == nil {
			fanRelation.LoyaltyPoints += loyaltyPoints
			fanRelation.TotalGiftAmount += int64(totalCost)
			fanRelation.LastGiftAt = ptrTimeNow()
			tx.Save(&fanRelation)
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, "failed to commit transaction")
		return
	}

	giftMsg := centrifugo.GiftMessage{
		Type:      "gift",
		Timestamp: time.Now().UnixMilli(),
	}
	giftMsg.Data.Sender.ID = userID
	giftMsg.Data.Sender.Nickname = user.Nickname
	if user.Nickname != "" {
		giftMsg.Data.Sender.Nickname = user.Nickname
	} else {
		giftMsg.Data.Sender.Nickname = user.Username
	}
	giftMsg.Data.Sender.Level = user.Level
	giftMsg.Data.Gift.ID = gift.ID
	giftMsg.Data.Gift.Name = gift.Name
	giftMsg.Data.Gift.Icon = gift.IconURL
	giftMsg.Data.Gift.Animation = gift.AnimationURL
	giftMsg.Data.Count = req.GiftCount
	giftMsg.Data.Combo = 1
	giftMsg.Data.TotalValue = totalCost

	client := centrifugo.NewClient("http://localhost:8000", "")
	channel := centrifugo.GetChannels(req.RoomID)[0]
	go client.Publish(channel, giftMsg)

	streamerRevenue := int64(0)
	if !isRelay {
		streamerRevenue = streamer.TotalRevenue
	}

	response.Success(c, gin.H{
		"transaction_id":    giftTx.ID,
		"gift_name":         gift.Name,
		"gift_count":        req.GiftCount,
		"total_cost":        totalCost,
		"remaining_balance": user.CoinBalance,
		"loyalty_points":    loyaltyPoints,
		"streamer_revenue":  streamerRevenue,
	})
}
