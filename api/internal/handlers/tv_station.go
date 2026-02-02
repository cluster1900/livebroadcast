package handlers

import (
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type PredefinedTVHandler struct{}

func NewPredefinedTVHandler() *PredefinedTVHandler {
	return &PredefinedTVHandler{}
}

type AddPredefinedTVRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" max=500"`
	SourceURL   string `json:"source_url" binding:"required"`
	Category    string `json:"category"`
	Country     string `json:"country"`
	Language    string `json:"language"`
	CoverURL    string `json:"cover_url"`
}

type TVStationResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SourceURL   string `json:"source_url"`
	Category    string `json:"category"`
	Country     string `json:"country"`
	Language    string `json:"language"`
	CoverURL    string `json:"cover_url"`
	IsActive    bool   `json:"is_active"`
}

func (h *PredefinedTVHandler) GetTVStations(c *gin.Context) {
	var stations []models.PredefinedRelay
	repository.DB.Where("is_active = ?", true).Order("sort_order ASC").Find(&stations)

	result := make([]TVStationResponse, 0, len(stations))
	for _, s := range stations {
		result = append(result, TVStationResponse{
			ID:          s.ID.String(),
			Name:        s.Name,
			Description: s.Description,
			SourceURL:   s.SourceURL,
			Category:    s.Category,
			Country:     s.Country,
			Language:    s.Language,
			CoverURL:    s.CoverURL,
			IsActive:    s.IsActive,
		})
	}

	response.Success(c, result)
}

func (h *PredefinedTVHandler) AddTVStation(c *gin.Context) {
	var req AddPredefinedTVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	station := models.PredefinedRelay{
		Name:        req.Name,
		Description: req.Description,
		SourceURL:   req.SourceURL,
		Category:    req.Category,
		Country:     req.Country,
		Language:    req.Language,
		CoverURL:    req.CoverURL,
		IsActive:    true,
	}

	if err := repository.DB.Create(&station).Error; err != nil {
		response.Fail(c, "failed to add TV station")
		return
	}

	response.Success(c, TVStationResponse{
		ID:          station.ID.String(),
		Name:        station.Name,
		Description: station.Description,
		SourceURL:   station.SourceURL,
		Category:    station.Category,
		Country:     station.Country,
		Language:    station.Language,
		CoverURL:    station.CoverURL,
		IsActive:    station.IsActive,
	})
}

func (h *PredefinedTVHandler) CreateRelaysFromTVStations(c *gin.Context) {
	var stations []models.PredefinedRelay
	repository.DB.Where("is_active = ?", true).Find(&stations)

	created := 0
	skipped := 0

	for _, s := range stations {
		var existingRelay models.RelayStream
		err := repository.DB.Where("source_url = ?", s.SourceURL).First(&existingRelay).Error
		if err == nil {
			skipped++
			continue
		}

		channelName := generateRelayChannelName()
		streamKey := generateRelayStreamKey()

		relay := models.RelayStream{
			Name:          s.Name,
			Description:   s.Description,
			SourceURL:     s.SourceURL,
			SourceType:    "rtmp",
			ChannelName:   channelName,
			StreamKey:     streamKey,
			RelayProtocol: "rtmp",
			Status:        "stopped",
			Category:      s.Category,
			CoverURL:      s.CoverURL,
			AutoStart:     true,
		}

		if err := repository.DB.Create(&relay).Error; err == nil {
			go startTVRelay(relay.SourceURL, relay.ChannelName)
			created++
		}
	}

	response.Success(c, gin.H{
		"message": "TV stations converted to relay streams",
		"created": created,
		"skipped": skipped,
	})
}

func startTVRelay(sourceURL, channelName string) {
	srsAPI := "http://localhost:1985/api/v1/relays"

	payload := `{
		"mode": "push",
		"source_url": "` + sourceURL + `",
		"destination": "rtmp://localhost/live/` + channelName + `"
	}`

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("POST", srsAPI, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		startTVFFmpegRelay(sourceURL, channelName)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		startTVFFmpegRelay(sourceURL, channelName)
	}
}

func startTVFFmpegRelay(sourceURL, channelName string) {
	cmd := "ffmpeg -i \"" + sourceURL + "\" -c copy -f flv \"rtmp://localhost/live/" + channelName + "\" -nostdin 2>&1 &"
	parts := strings.Split(cmd, " ")
	head := parts[0]
	parts = parts[1:]

	execCmd := strings.Trim(head, "`\"'")
	arg := "-c"
	fullCmd := strings.Join(parts, " ")

	exec.Command(execCmd, arg, fullCmd).Start()
}
