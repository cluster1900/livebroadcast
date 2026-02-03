package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type HistoryHandler struct{}

func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{}
}

func (h *HistoryHandler) GetWatchHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	userUUID := uuid.MustParse(userID)
	var history []struct {
		ID            string `json:"id"`
		RoomID        string `json:"room_id"`
		RoomTitle     string `json:"room_title"`
		StreamerID    string `json:"streamer_id"`
		StreamerName  string `json:"streamer_name"`
		CoverURL      string `json:"cover_url"`
		WatchDuration int    `json:"watch_duration"`
		StartTime     string `json:"start_time"`
	}

	repository.DB.Raw(`
		SELECT 
			h.id, h.room_id, r.title as room_title,
			r.streamer_id, u.nickname as streamer_name,
			r.cover_url, h.watch_duration, h.created_at as start_time
		FROM watch_histories h
		JOIN live_rooms r ON r.id = h.room_id
		JOIN users u ON u.id = r.streamer_id
		WHERE h.user_id = ?
		ORDER BY h.created_at DESC
		LIMIT 50
	`, userUUID).Scan(&history)

	response.Success(c, history)
}

func (h *HistoryHandler) AddWatchHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		RoomID        string `json:"room_id" binding:"required"`
		WatchDuration int    `json:"watch_duration"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userUUID := uuid.MustParse(userID)
	roomUUID := uuid.MustParse(req.RoomID)

	var history models.WatchHistory
	if err := repository.DB.Where("user_id = ? AND room_id = ?", userUUID, roomUUID).First(&history).Error; err == nil {
		history.WatchDuration += req.WatchDuration
		repository.DB.Save(&history)
	} else {
		history = models.WatchHistory{
			ID:            uuid.New(),
			UserID:        userUUID,
			RoomID:        roomUUID,
			WatchDuration: req.WatchDuration,
		}
		repository.DB.Create(&history)
	}

	response.Success(c, gin.H{"message": "记录成功"})
}

func (h *HistoryHandler) ClearWatchHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	userUUID := uuid.MustParse(userID)
	repository.DB.Where("user_id = ?", userUUID).Delete(&models.WatchHistory{})

	response.Success(c, gin.H{"message": "清除成功"})
}

func (h *HistoryHandler) DeleteWatchHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	historyID := c.Param("id")
	userUUID := uuid.MustParse(userID)

	if err := repository.DB.Where("id = ? AND user_id = ?", historyID, userUUID).Delete(&models.WatchHistory{}).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
