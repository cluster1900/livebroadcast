package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ok",
		"message": "server is running",
	})
}
