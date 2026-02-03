package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type GiftInventoryHandler struct{}

func NewGiftInventoryHandler() *GiftInventoryHandler {
	return &GiftInventoryHandler{}
}

func (h *GiftInventoryHandler) GetInventory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	userUUID := uuid.MustParse(userID)
	var inventory []struct {
		GiftID    int    `json:"gift_id"`
		Name      string `json:"name"`
		IconURL   string `json:"icon_url"`
		Count     int    `json:"count"`
		CoinPrice int    `json:"coin_price"`
		Category  string `json:"category"`
	}

	repository.DB.Raw(`
		SELECT g.id as gift_id, g.name, g.icon_url, COALESCE(i.count, 0) as count, g.coin_price, g.category
		FROM gifts g
		LEFT JOIN gift_inventories i ON i.gift_id = g.id AND i.user_id = ?
		WHERE g.is_active = true AND (i.count > 0 OR g.category = 'special')
		ORDER BY g.sort_order
	`, userUUID).Scan(&inventory)

	response.Success(c, inventory)
}

func (h *GiftInventoryHandler) AddGift(userID uuid.UUID, giftID int, count int) error {
	var inventory models.GiftInventory
	if err := repository.DB.Where("user_id = ? AND gift_id = ?", userID, giftID).First(&inventory).Error; err == nil {
		inventory.Count += count
		return repository.DB.Save(&inventory).Error
	}

	inventory = models.GiftInventory{
		UserID: userID,
		GiftID: giftID,
		Count:  count,
	}
	return repository.DB.Create(&inventory).Error
}

func (h *GiftInventoryHandler) UseGift(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		GiftID int `json:"gift_id" binding:"required"`
		Count  int `json:"count" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userUUID := uuid.MustParse(userID)
	var inventory models.GiftInventory
	if err := repository.DB.Where("user_id = ? AND gift_id = ?", userUUID, req.GiftID).First(&inventory).Error; err != nil {
		response.Fail(c, "礼物不存在或数量不足")
		return
	}

	if inventory.Count < req.Count {
		response.Fail(c, "礼物数量不足")
		return
	}

	inventory.Count -= req.Count
	repository.DB.Save(&inventory)

	response.Success(c, gin.H{"message": "使用成功", "remaining": inventory.Count})
}
