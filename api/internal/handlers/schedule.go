package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type ScheduleHandler struct{}

func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{}
}

type CreateScheduleRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description,max=1000"`
	Category    string `json:"category"`
	CoverURL    string `json:"cover_url"`
	StartTime   string `json:"start_time" binding:"required"`
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		response.BadRequest(c, "时间格式错误，请使用RFC3339格式")
		return
	}

	if startTime.Before(time.Now()) {
		response.BadRequest(c, "开始时间不能早于当前时间")
		return
	}

	schedule := models.LiveSchedule{
		ID:          uuid.New(),
		StreamerID:  uuid.MustParse(userID),
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		CoverURL:    req.CoverURL,
		StartTime:   startTime,
		Status:      "scheduled",
	}

	if err := repository.DB.Create(&schedule).Error; err != nil {
		response.Fail(c, "创建失败")
		return
	}

	response.Success(c, gin.H{"message": "创建成功", "data": schedule})
}

func (h *ScheduleHandler) GetMySchedules(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	userUUID := uuid.MustParse(userID)
	var schedules []models.LiveSchedule
	repository.DB.Where("streamer_id = ?", userUUID).Order("start_time ASC").Find(&schedules)

	response.Success(c, schedules)
}

func (h *ScheduleHandler) GetUpcomingSchedules(c *gin.Context) {
	var schedules []struct {
		ID           string `json:"id"`
		StreamerID   string `json:"streamer_id"`
		StreamerName string `json:"streamer_name"`
		Title        string `json:"title"`
		Category     string `json:"category"`
		CoverURL     string `json:"cover_url"`
		StartTime    string `json:"start_time"`
	}

	repository.DB.Raw(`
		SELECT s.id, s.streamer_id, u.nickname as streamer_name, 
		       s.title, s.category, s.cover_url, s.start_time
		FROM live_schedules s
		JOIN users u ON u.id = s.streamer_id
		WHERE s.start_time > NOW() AND s.status = 'scheduled'
		ORDER BY s.start_time ASC
		LIMIT 50
	`).Scan(&schedules)

	response.Success(c, schedules)
}

func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	scheduleID := c.Param("id")
	userUUID := uuid.MustParse(userID)

	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	startTime, _ := time.Parse(time.RFC3339, req.StartTime)

	updates := map[string]interface{}{
		"title":       req.Title,
		"description": req.Description,
		"category":    req.Category,
		"cover_url":   req.CoverURL,
		"start_time":  startTime,
	}

	if err := repository.DB.Model(&models.LiveSchedule{}).Where("id = ? AND streamer_id = ?", scheduleID, userUUID).Updates(updates).Error; err != nil {
		response.Fail(c, "更新失败")
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

func (h *ScheduleHandler) CancelSchedule(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	scheduleID := c.Param("id")
	userUUID := uuid.MustParse(userID)

	if err := repository.DB.Model(&models.LiveSchedule{}).Where("id = ? AND streamer_id = ?", scheduleID, userUUID).Update("status", "cancelled").Error; err != nil {
		response.Fail(c, "取消失败")
		return
	}

	response.Success(c, gin.H{"message": "已取消"})
}

func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	scheduleID := c.Param("id")
	userUUID := uuid.MustParse(userID)

	if err := repository.DB.Where("id = ? AND streamer_id = ?", scheduleID, userUUID).Delete(&models.LiveSchedule{}).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
