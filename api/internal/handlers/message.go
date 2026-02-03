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

type MessageHandler struct {
	centrifugoClient *centrifugo.Client
}

func NewMessageHandler(centrifugoClient *centrifugo.Client) *MessageHandler {
	return &MessageHandler{centrifugoClient: centrifugoClient}
}

type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required,max=500"`
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	senderID := c.GetString("user_id")
	if senderID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if senderID == req.ReceiverID {
		response.BadRequest(c, "不能给自己发送私信")
		return
	}

	receiverUUID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		response.BadRequest(c, "接收者ID无效")
		return
	}

	var receiver models.User
	if err := repository.DB.First(&receiver, "id = ?", req.ReceiverID).Error; err != nil {
		response.BadRequest(c, "接收者不存在")
		return
	}

	senderUUID := uuid.MustParse(senderID)
	message := models.PrivateMessage{
		ID:         uuid.New(),
		SenderID:   senderUUID,
		ReceiverID: receiverUUID,
		Content:    req.Content,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}

	if err := repository.DB.Create(&message).Error; err != nil {
		response.Fail(c, "发送失败")
		return
	}

	if h.centrifugoClient != nil {
		h.centrifugoClient.Publish("user:"+req.ReceiverID, map[string]interface{}{
			"type":        "private_message",
			"id":          message.ID.String(),
			"sender_id":   senderID,
			"sender_name": getNicknameByID(senderUUID),
			"content":     req.Content,
			"created_at":  message.CreatedAt.Format(time.RFC3339),
		})
	}

	response.Success(c, gin.H{"message": "发送成功", "data": message})
}

func getNicknameByID(userID uuid.UUID) string {
	var user models.User
	repository.DB.Select("nickname").First(&user, "id = ?", userID)
	if user.Nickname != "" {
		return user.Nickname
	}
	return "用户"
}

func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	userUUID := uuid.MustParse(userID)
	var conversations []struct {
		UserID      string    `json:"user_id"`
		Nickname    string    `json:"nickname"`
		AvatarURL   string    `json:"avatar_url"`
		LastMessage string    `json:"last_message"`
		LastTime    time.Time `json:"last_time"`
		UnreadCount int       `json:"unread_count"`
	}

	repository.DB.Raw(`
		SELECT 
			CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as user_id,
			u.nickname, u.avatar_url,
			m.content as last_message,
			m.created_at as last_time,
			(SELECT COUNT(*) FROM private_messages WHERE ((sender_id = ? AND receiver_id = user_id) OR (sender_id = user_id AND receiver_id = ?)) AND is_read = FALSE AND sender_id != ?) as unread_count
		FROM private_messages m
		JOIN users u ON u.id = CASE WHEN m.sender_id = ? THEN m.receiver_id ELSE m.sender_id END
		WHERE ? IN (m.sender_id, m.receiver_id)
		GROUP BY user_id, u.nickname, u.avatar_url, m.content, m.created_at
		ORDER BY last_time DESC
	`, userUUID, userUUID, userUUID, userUUID, userUUID, userUUID).Scan(&conversations)

	response.Success(c, conversations)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	otherID := c.Param("user_id")
	otherUUID, err := uuid.Parse(otherID)
	if err != nil {
		response.BadRequest(c, "用户ID无效")
		return
	}

	var messages []models.PrivateMessage
	userUUID := uuid.MustParse(userID)

	if err := repository.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userUUID, otherUUID, otherUUID, userUUID,
	).Order("created_at DESC").Limit(100).Find(&messages).Error; err != nil {
		response.Fail(c, "获取消息失败")
		return
	}

	repository.DB.Model(&models.PrivateMessage{}).Where(
		"sender_id = ? AND receiver_id = ? AND is_read = ?",
		otherUUID, userUUID, false,
	).Update("is_read", true)

	response.Success(c, messages)
}

func (h *MessageHandler) GetUnreadMessageCount(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var count int64
	userUUID := uuid.MustParse(userID)
	repository.DB.Model(&models.PrivateMessage{}).Where("receiver_id = ? AND is_read = ?", userUUID, false).Count(&count)

	response.Success(c, gin.H{"count": int(count)})
}

func (h *MessageHandler) DeleteConversation(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	otherID := c.Param("user_id")
	otherUUID := uuid.MustParse(otherID)
	userUUID := uuid.MustParse(userID)

	if err := repository.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userUUID, otherUUID, otherUUID, userUUID,
	).Delete(&models.PrivateMessage{}).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
