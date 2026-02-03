package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var notifications []models.Notification
	if err := repository.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&notifications).Error; err != nil {
		response.Fail(c, "获取通知失败")
		return
	}

	response.Success(c, notifications)
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var count int64
	repository.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count)

	response.Success(c, gin.H{"count": count})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	notificationID := c.Param("id")
	var notification models.Notification
	if err := repository.DB.Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		response.Fail(c, "通知不存在")
		return
	}

	notification.IsRead = true
	repository.DB.Save(&notification)

	response.Success(c, gin.H{"message": "已标记为已读"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	repository.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true)

	response.Success(c, gin.H{"message": "已全部标记为已读"})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	notificationID := c.Param("id")
	if err := repository.DB.Where("id = ? AND user_id = ?", notificationID, userID).Delete(&models.Notification{}).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

func CreateNotification(userID uuid.UUID, notifType, title, content, link string) error {
	notification := models.Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Content:   content,
		Link:      link,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	return repository.DB.Create(&notification).Error
}
