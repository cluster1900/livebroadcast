package handlers

import (
	"github.com/gin-gonic/gin"
)

type CentrifugoHandler struct{}

func NewCentrifugoHandler() *CentrifugoHandler {
	return &CentrifugoHandler{}
}

type CentrifugoTokenResponse struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

func (h *CentrifugoHandler) GetToken(c *gin.Context) {
	userID := c.GetString("user_id")

	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{
		"token": "centrifugo_token_" + userID,
		"url":   "ws://localhost:8000/connection/websocket",
	}})
}
