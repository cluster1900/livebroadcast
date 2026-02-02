package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type RelayHandler struct{}

func NewRelayHandler() *RelayHandler {
	return &RelayHandler{}
}

type CreateRelayRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" max=500"`
	SourceURL   string `json:"source_url" binding:"required"`
	SourceType  string `json:"source_type"`
	Category    string `json:"category"`
	CoverURL    string `json:"cover_url"`
	AutoStart   bool   `json:"auto_start"`
}

type UpdateRelayRequest struct {
	Name        string `json:"name" max=100"`
	Description string `json:"description" max=500"`
	Category    string `json:"category"`
	CoverURL    string `json:"cover_url"`
	AutoStart   *bool  `json:"auto_start"`
}

type RelayResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SourceURL   string `json:"source_url"`
	ChannelName string `json:"channel_name"`
	StreamURL   string `json:"stream_url"`
	Status      string `json:"status"`
	Category    string `json:"category"`
	CoverURL    string `json:"cover_url"`
	ViewCount   int64  `json:"view_count"`
	PeakOnline  int64  `json:"peak_online"`
	IsActive    bool   `json:"is_active"`
}

func generateRelayStreamKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateRelayChannelName() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "relay_" + hex.EncodeToString(b)
}

func (h *RelayHandler) CreateRelay(c *gin.Context) {
	var req CreateRelayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	sourceType := req.SourceType
	if sourceType == "" {
		sourceType = "rtmp"
	}

	if !isValidRelaySourceType(sourceType) {
		response.BadRequest(c, "invalid source_type, supported: rtmp, http-flv, hls")
		return
	}

	channelName := generateRelayChannelName()
	streamKey := generateRelayStreamKey()

	relay := models.RelayStream{
		ID:            uuid.New(),
		Name:          req.Name,
		Description:   req.Description,
		SourceURL:     req.SourceURL,
		SourceType:    sourceType,
		ChannelName:   channelName,
		StreamKey:     streamKey,
		RelayProtocol: "rtmp",
		Status:        "stopped",
		Category:      req.Category,
		CoverURL:      req.CoverURL,
		AutoStart:     req.AutoStart,
	}

	if err := repository.DB.Create(&relay).Error; err != nil {
		response.Fail(c, "failed to create relay stream")
		return
	}

	if req.AutoStart {
		go h.startRelay(relay.ID.String())
	}

	response.Success(c, RelayResponse{
		ID:          relay.ID.String(),
		Name:        relay.Name,
		Description: relay.Description,
		SourceURL:   relay.SourceURL,
		ChannelName: relay.ChannelName,
		StreamURL:   fmt.Sprintf("rtmp://localhost/live/%s", relay.ChannelName),
		Status:      relay.Status,
		Category:    relay.Category,
		CoverURL:    relay.CoverURL,
		ViewCount:   relay.ViewCount,
		PeakOnline:  relay.PeakOnline,
		IsActive:    relay.AutoStart,
	})
}

func (h *RelayHandler) GetRelays(c *gin.Context) {
	var relays []models.RelayStream
	repository.DB.Order("created_at DESC").Find(&relays)

	result := make([]RelayResponse, 0, len(relays))
	for _, r := range relays {
		result = append(result, RelayResponse{
			ID:          r.ID.String(),
			Name:        r.Name,
			Description: r.Description,
			SourceURL:   maskRelaySourceURL(r.SourceURL),
			ChannelName: r.ChannelName,
			StreamURL:   fmt.Sprintf("rtmp://localhost/live/%s", r.ChannelName),
			Status:      r.Status,
			Category:    r.Category,
			CoverURL:    r.CoverURL,
			ViewCount:   r.ViewCount,
			PeakOnline:  r.PeakOnline,
			IsActive:    r.Status == "running",
		})
	}

	response.Success(c, result)
}

func (h *RelayHandler) GetRelay(c *gin.Context) {
	id := c.Param("id")

	var relay models.RelayStream
	if err := repository.DB.First(&relay, "id = ?", id).Error; err != nil {
		response.BadRequest(c, "relay stream not found")
		return
	}

	response.Success(c, RelayResponse{
		ID:          relay.ID.String(),
		Name:        relay.Name,
		Description: relay.Description,
		SourceURL:   relay.SourceURL,
		ChannelName: relay.ChannelName,
		StreamURL:   fmt.Sprintf("rtmp://localhost/live/%s", relay.ChannelName),
		Status:      relay.Status,
		Category:    relay.Category,
		CoverURL:    relay.CoverURL,
		ViewCount:   relay.ViewCount,
		PeakOnline:  relay.PeakOnline,
		IsActive:    relay.Status == "running",
	})
}

func (h *RelayHandler) UpdateRelay(c *gin.Context) {
	id := c.Param("id")

	var relay models.RelayStream
	if err := repository.DB.First(&relay, "id = ?", id).Error; err != nil {
		response.BadRequest(c, "relay stream not found")
		return
	}

	var req UpdateRelayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.CoverURL != "" {
		updates["cover_url"] = req.CoverURL
	}
	if req.AutoStart != nil {
		updates["auto_start"] = *req.AutoStart
	}

	if len(updates) > 0 {
		repository.DB.Model(&relay).Updates(updates)
	}

	response.Success(c, gin.H{
		"message": "relay stream updated",
		"id":      relay.ID.String(),
	})
}

func (h *RelayHandler) DeleteRelay(c *gin.Context) {
	id := c.Param("id")

	var relay models.RelayStream
	if err := repository.DB.First(&relay, "id = ?", id).Error; err != nil {
		response.BadRequest(c, "relay stream not found")
		return
	}

	if relay.Status == "running" {
		h.stopRelay(id)
	}

	repository.DB.Delete(&relay)

	response.Success(c, gin.H{
		"message": "relay stream deleted",
		"id":      id,
	})
}

func (h *RelayHandler) StartRelay(c *gin.Context) {
	id := c.Param("id")

	var relay models.RelayStream
	if err := repository.DB.First(&relay, "id = ?", id).Error; err != nil {
		response.BadRequest(c, "relay stream not found")
		return
	}

	if relay.Status == "running" {
		response.BadRequest(c, "relay stream is already running")
		return
	}

	go h.startRelay(id)

	response.Success(c, gin.H{
		"message":    "relay stream starting",
		"id":         id,
		"stream_url": fmt.Sprintf("rtmp://localhost/live/%s", relay.ChannelName),
	})
}

func (h *RelayHandler) StopRelay(c *gin.Context) {
	id := c.Param("id")

	var relay models.RelayStream
	if err := repository.DB.First(&relay, "id = ?", id).Error; err != nil {
		response.BadRequest(c, "relay stream not found")
		return
	}

	if relay.Status != "running" {
		response.BadRequest(c, "relay stream is not running")
		return
	}

	h.stopRelay(id)

	response.Success(c, gin.H{
		"message": "relay stream stopped",
		"id":      id,
	})
}

func (h *RelayHandler) startRelay(id string) {
	var relay models.RelayStream
	if err := repository.DB.First(&relay, "id = ?", id).Error; err != nil {
		return
	}

	repository.DB.Model(&relay).Update("status", "starting")

	err := startRelayProcess(relay.SourceURL, relay.ChannelName)
	if err != nil {
		repository.DB.Model(&relay).Update("status", "error")
		h.logRelayEvent(relay.ID.String(), "start_error", "", err.Error())
		return
	}

	repository.DB.Model(&relay).Update("status", "running")
	h.logRelayEvent(relay.ID.String(), "started", "", "")
}

func (h *RelayHandler) stopRelay(id string) {
	var relay models.RelayStream
	if err := repository.DB.First(&relay, "id = ?", id).Error; err != nil {
		return
	}

	stopRelayProcess(relay.ChannelName)

	repository.DB.Model(&relay).Update("status", "stopped")
	h.logRelayEvent(relay.ID.String(), "stopped", "", "")
}

func (h *RelayHandler) logRelayEvent(relayID, eventType, eventData, errorMsg string) {
	log := models.RelayStreamLog{
		RelayStreamID: uuid.MustParse(relayID),
		EventType:     eventType,
		EventData:     eventData,
		ErrorMessage:  errorMsg,
	}
	repository.DB.Create(&log)
}

func startRelayProcess(sourceURL, channelName string) error {
	srsAPI := "http://localhost:1985/api/v1/relays"

	payload := fmt.Sprintf(`{
		"mode": "push",
		"source_url": "%s",
		"destination": "rtmp://localhost/live/%s"
	}`, sourceURL, channelName)

	req, _ := http.NewRequest("POST", srsAPI, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return startFFmpegRelay(sourceURL, channelName)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(body), "already") || strings.Contains(string(body), "exists") {
		return nil
	}

	return startFFmpegRelay(sourceURL, channelName)
}

func stopRelayProcess(channelName string) {
	srsAPI := fmt.Sprintf("http://localhost:1985/api/v1/relays/live%%2F%s", channelName)

	req, _ := http.NewRequest("DELETE", srsAPI, nil)
	client := &http.Client{Timeout: 10 * time.Second}
	client.Do(req)

	stopFFmpegRelay(channelName)
}

func startFFmpegRelay(sourceURL, channelName string) error {
	log.Printf("Starting FFmpeg relay: %s -> rtmp://localhost/live/%s", sourceURL, channelName)

	cmd := exec.Command("ffmpeg",
		"-re",
		"-i", sourceURL,
		"-c", "copy",
		"-f", "flv",
		fmt.Sprintf("rtmp://localhost/live/%s", channelName),
		"-nostdin",
	)

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start FFmpeg: %v", err)
		return err
	}

	log.Printf("FFmpeg relay started with PID: %d", cmd.Process.Pid)
	return nil
}

func stopFFmpegRelay(channelName string) {
	log.Printf("Stopping FFmpeg relay for: %s", channelName)

	cmd := exec.Command("pkill", "-f", fmt.Sprintf("ffmpeg.*%s", channelName))
	cmd.Run()
}

func isValidRelaySourceType(t string) bool {
	validTypes := []string{"rtmp", "http-flv", "hls", "rtsp"}
	for _, vt := range validTypes {
		if t == vt {
			return true
		}
	}
	return false
}

func maskRelaySourceURL(url string) string {
	if len(url) > 50 {
		return url[:50] + "..."
	}
	return url
}
