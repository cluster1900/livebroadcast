package handlers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/centrifugo"
	"github.com/huya_live/api/pkg/response"
)

type DanmuHandler struct{}

func NewDanmuHandler() *DanmuHandler {
	return &DanmuHandler{}
}

type SendDanmuRequest struct {
	RoomID  string `json:"room_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1,max=100"`
	Color   string `json:"color"`
}

func (h *DanmuHandler) SendDanmu(c *gin.Context) {
	userID := c.GetString("user_id")
	var req SendDanmuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	var room models.LiveRoom
	var relay models.RelayStream
	roomFound := false

	if err := repository.DB.Where("id = ? AND status = ?", req.RoomID, "live").First(&room).Error; err == nil {
		roomFound = true
	} else {
		if err := repository.DB.Where("id = ? AND status = ?", req.RoomID, "running").First(&relay).Error; err == nil {
			roomFound = true
		}
	}

	if !roomFound {
		response.BadRequest(c, "room not found or not live")
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		response.BadRequest(c, "content cannot be empty")
		return
	}

	sensitiveWords := []models.SensitiveWord{}
	repository.DB.Where("is_active = ?", true).Find(&sensitiveWords)
	for _, word := range sensitiveWords {
		if strings.Contains(strings.ToLower(content), strings.ToLower(word.Word)) {
			response.BadRequest(c, "content contains sensitive words")
			return
		}
	}

	var user models.User
	if err := repository.DB.First(&user, "id = ?", userID).Error; err != nil {
		response.BadRequest(c, "user not found")
		return
	}

	danmuColor := req.Color
	if danmuColor == "" {
		danmuColor = "#FFFFFF"
	}

	danmuMsg := centrifugo.DanmuMessage{
		Type:      "danmu",
		Timestamp: time.Now().UnixMilli(),
	}
	danmuMsg.Data.ID = uuid.New().String()
	danmuMsg.Data.UserID = userID
	danmuMsg.Data.Nickname = user.Nickname
	if user.Nickname != "" {
		danmuMsg.Data.Nickname = user.Nickname
	} else {
		danmuMsg.Data.Nickname = user.Username
	}
	danmuMsg.Data.Level = user.Level
	danmuMsg.Data.Avatar = user.AvatarURL
	danmuMsg.Data.Content = content
	danmuMsg.Data.Color = danmuColor

	client := centrifugo.NewClient("http://localhost:8000", "")
	channel := centrifugo.GetChannels(req.RoomID)[0]
	if err := client.Publish(channel, danmuMsg); err != nil {
		response.Fail(c, "failed to send danmu: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message_id": danmuMsg.Data.ID,
		"content":    content,
		"timestamp":  danmuMsg.Timestamp,
	})
}
