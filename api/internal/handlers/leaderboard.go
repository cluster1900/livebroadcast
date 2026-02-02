package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type LeaderboardHandler struct{}

func NewLeaderboardHandler() *LeaderboardHandler {
	return &LeaderboardHandler{}
}

type LeaderboardEntry struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Score      int64  `json:"score"`
	Rank       int    `json:"rank"`
	Level      int    `json:"level"`
	StreamerID string `json:"streamer_id"`
}

func (h *LeaderboardHandler) GetRoomLeaderboard(c *gin.Context) {
	roomID := c.Param("room_id")

	var transactions []models.GiftTransaction
	repository.DB.Where("room_id = ?", roomID).
		Select("receiver_id, SUM(coin_amount) as total").
		Group("receiver_id").
		Order("total DESC").
		Limit(10).
		Find(&transactions)

	result := make([]LeaderboardEntry, 0, len(transactions))
	for i, tx := range transactions {
		var streamer models.Streamer
		repository.DB.Where("user_id = ?", tx.ReceiverID).First(&streamer)

		var user models.User
		repository.DB.Where("id = ?", tx.ReceiverID).First(&user)

		name := user.Nickname
		if name == "" {
			name = user.Username
		}

		result = append(result, LeaderboardEntry{
			ID:         user.ID.String(),
			Name:       name,
			Avatar:     user.AvatarURL,
			Score:      int64(tx.CoinAmount),
			Rank:       i + 1,
			Level:      user.Level,
			StreamerID: streamer.UserID.String(),
		})
	}

	response.Success(c, result)
}

func (h *LeaderboardHandler) GetGlobalLeaderboard(c *gin.Context) {
	limit := 20
	streamers := []models.Streamer{}
	repository.DB.Order("total_revenue DESC").Limit(limit).Find(&streamers)

	result := make([]LeaderboardEntry, 0, len(streamers))
	for i, streamer := range streamers {
		var user models.User
		repository.DB.Where("id = ?", streamer.UserID).First(&user)

		name := user.Nickname
		if name == "" {
			name = user.Username
		}

		result = append(result, LeaderboardEntry{
			ID:         user.ID.String(),
			Name:       name,
			Avatar:     user.AvatarURL,
			Score:      streamer.TotalRevenue,
			Rank:       i + 1,
			Level:      user.Level,
			StreamerID: streamer.UserID.String(),
		})
	}

	response.Success(c, result)
}

func (h *LeaderboardHandler) GetRichList(c *gin.Context) {
	limit := 20
	users := []models.User{}
	repository.DB.Order("coin_balance DESC").Limit(limit).Find(&users)

	result := make([]LeaderboardEntry, 0, len(users))
	for i, user := range users {
		name := user.Nickname
		if name == "" {
			name = user.Username
		}

		result = append(result, LeaderboardEntry{
			ID:     user.ID.String(),
			Name:   name,
			Avatar: user.AvatarURL,
			Score:  int64(user.CoinBalance),
			Rank:   i + 1,
			Level:  user.Level,
		})
	}

	response.Success(c, result)
}

func (h *LeaderboardHandler) GetCategories(c *gin.Context) {
	categories := []string{
		"娱乐",
		"游戏",
		"美食",
		"音乐",
		"舞蹈",
		"户外",
		"科技",
		"体育",
		"汽车",
		"时尚",
		"教育",
		"财经",
		"新闻",
		"综合",
		"测试",
	}

	result := make([]gin.H, 0, len(categories))
	for i, cat := range categories {
		result = append(result, gin.H{
			"id":    i + 1,
			"name":  cat,
			"icon":  "https://placeholder.com/category-" + cat + ".png",
			"count": 0,
		})
	}

	response.Success(c, result)
}

func (h *LeaderboardHandler) GetOnlineCount(c *gin.Context) {
	var liveCount int64
	repository.DB.Model(&models.LiveRoom{}).Where("status = ?", "live").Count(&liveCount)

	var relayCount int64
	repository.DB.Model(&models.RelayStream{}).Where("status = ?", "running").Count(&relayCount)

	response.Success(c, gin.H{
		"live_rooms":    liveCount,
		"relay_streams": relayCount,
		"total":         liveCount + relayCount,
	})
}
