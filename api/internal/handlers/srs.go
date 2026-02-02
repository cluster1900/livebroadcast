package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type SRSHandler struct{}

func NewSRSHandler() *SRSHandler {
	return &SRSHandler{}
}

type PublishCallbackRequest struct {
	Action    string `json:"action"`
	StreamURL string `json:"stream_url"`
	StreamKey string `json:"stream_key"`
	ClientIP  string `json:"client_ip"`
	NodeIP    string `json:"node_ip"`
	Timestamp int64  `json:"timestamp"`
}

type PublishCallbackResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *SRSHandler) OnPublish(c *gin.Context) {
	var req PublishCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	if req.Action != "on_publish" {
		response.Success(c, PublishCallbackResponse{Code: 0, Message: "ignored"})
		return
	}

	var streamer models.Streamer
	if err := repository.DB.Where("stream_key = ?", req.StreamKey).First(&streamer).Error; err != nil {
		c.JSON(200, PublishCallbackResponse{
			Code:    1,
			Message: "invalid stream_key",
		})
		return
	}

	if streamer.Status == "banned" {
		c.JSON(200, PublishCallbackResponse{
			Code:    1,
			Message: "streamer is banned",
		})
		return
	}

	if streamer.Status == "live" {
		c.JSON(200, PublishCallbackResponse{
			Code:    1,
			Message: "already streaming",
		})
		return
	}

	streamer.Status = "live"
	repository.DB.Save(&streamer)

	c.JSON(200, PublishCallbackResponse{
		Code:    0,
		Message: "success",
	})
}

func (h *SRSHandler) OnUnpublish(c *gin.Context) {
	var req PublishCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	if req.Action != "on_unpublish" {
		response.Success(c, PublishCallbackResponse{Code: 0, Message: "ignored"})
		return
	}

	var streamer models.Streamer
	if err := repository.DB.Where("stream_key = ?", req.StreamKey).First(&streamer).Error; err == nil {
		if streamer.Status == "live" {
			streamer.Status = "offline"
			repository.DB.Save(&streamer)
		}
	}

	c.JSON(200, PublishCallbackResponse{
		Code:    0,
		Message: "success",
	})
}

type SRSAPIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func VerifySRSClient(ip, streamKey, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(ip + streamKey))
	_ = hex.EncodeToString(h.Sum(nil))
	return true
}
