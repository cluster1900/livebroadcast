package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type LiveHandler struct{}

func NewLiveHandler() *LiveHandler {
	return &LiveHandler{}
}

type CreateRoomRequest struct {
	Title    string `json:"title" binding:"required,min=2,max=200"`
	Category string `json:"category"`
	CoverURL string `json:"cover_url"`
}

func (h *LiveHandler) CreateRoom(c *gin.Context) {
	userID := c.GetString("user_id")
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	var streamer models.Streamer
	if err := repository.DB.Where("user_id = ?", userID).First(&streamer).Error; err != nil {
		response.BadRequest(c, "you are not a streamer")
		return
	}

	existingRoom := models.LiveRoom{}
	if err := repository.DB.Where("streamer_id = ? AND status = ?", userID, "live").First(&existingRoom).Error; err == nil {
		response.BadRequest(c, "you already have a live room")
		return
	}

	channelName := "live_" + uuid.New().String()[:8]

	room := models.LiveRoom{
		StreamerID:  uuid.MustParse(userID),
		Title:       req.Title,
		Category:    req.Category,
		CoverURL:    req.CoverURL,
		ChannelName: channelName,
		Status:      "live",
		StartAt:     ptrTimeNow(),
	}

	if err := repository.DB.Create(&room).Error; err != nil {
		response.Fail(c, "failed to create room")
		return
	}

	response.Success(c, gin.H{
		"room_id":      room.ID.String(),
		"channel_name": channelName,
		"stream_url":   "rtmp://localhost/live/" + channelName,
		"status":       "live",
	})
}

func (h *LiveHandler) EndRoom(c *gin.Context) {
	userID := c.GetString("user_id")
	roomID := c.Param("id")

	var room models.LiveRoom
	if err := repository.DB.Where("id = ? AND streamer_id = ?", roomID, userID).First(&room).Error; err != nil {
		response.BadRequest(c, "room not found")
		return
	}

	if room.Status != "live" {
		response.BadRequest(c, "room is not live")
		return
	}

	now := ptrTimeNow()
	room.Status = "ended"
	room.EndAt = now

	if err := repository.DB.Save(&room).Error; err != nil {
		response.Fail(c, "failed to end room")
		return
	}

	response.Success(c, gin.H{
		"message": "room ended successfully",
		"room_id": room.ID.String(),
	})
}

type UpdateRoomRequest struct {
	Title    string `json:"title" max=200"`
	Category string `json:"category"`
	CoverURL string `json:"cover_url"`
}

func (h *LiveHandler) UpdateRoom(c *gin.Context) {
	userID := c.GetString("user_id")
	roomID := c.Param("id")

	var room models.LiveRoom
	if err := repository.DB.Where("id = ? AND streamer_id = ?", roomID, userID).First(&room).Error; err != nil {
		response.BadRequest(c, "room not found")
		return
	}

	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.CoverURL != "" {
		updates["cover_url"] = req.CoverURL
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		repository.DB.Model(&room).Updates(updates)
	}

	response.Success(c, gin.H{
		"message": "room updated successfully",
		"room_id": room.ID.String(),
	})
}

type GetRoomResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Category     string `json:"category"`
	CoverURL     string `json:"cover_url"`
	ChannelName  string `json:"channel_name"`
	Status       string `json:"status"`
	StreamerID   string `json:"streamer_id"`
	StreamerName string `json:"streamer_name"`
	StartAt      string `json:"start_at"`
	PeakOnline   int    `json:"peak_online"`
	TotalViews   int    `json:"total_views"`
	StreamURL    string `json:"stream_url,omitempty"`
	FLVURL       string `json:"flv_url,omitempty"`
	HLSURL       string `json:"hls_url,omitempty"`
}

func (h *LiveHandler) GetRoom(c *gin.Context) {
	roomID := c.Param("id")

	var room models.LiveRoom
	if err := repository.DB.First(&room, "id = ?", roomID).Error; err != nil {
		var relay models.RelayStream
		if err := repository.DB.First(&relay, "id = ?", roomID).Error; err != nil {
			response.BadRequest(c, "room not found")
			return
		}

		resp := GetRoomResponse{
			ID:           relay.ID.String(),
			Title:        relay.Name,
			Category:     relay.Category,
			CoverURL:     relay.CoverURL,
			ChannelName:  relay.ChannelName,
			Status:       relay.Status,
			StreamerID:   "relay-" + relay.ID.String()[:8],
			StreamerName: relay.Name,
			StartAt:      formatTimeFromTime(relay.CreatedAt),
			PeakOnline:   int(relay.PeakOnline),
			TotalViews:   int(relay.ViewCount),
		}

		if relay.Status == "running" {
			// Always return SRS URLs for relay streams
			// The actual relay is handled by FFmpeg pushing to SRS
			resp.StreamURL = "rtmp://localhost/live/" + relay.ChannelName
			resp.FLVURL = "http://localhost:8080/live/" + relay.ChannelName + ".flv"
			resp.HLSURL = "http://localhost:8080/live/" + relay.ChannelName + ".m3u8"
		}

		response.Success(c, resp)
		return
	}

	streamerName := ""
	if room.StreamerID != uuid.Nil {
		var user models.User
		if err := repository.DB.Select("nickname").Where("id = ?", room.StreamerID).First(&user).Error; err == nil {
			streamerName = user.Nickname
			if streamerName == "" {
				streamerName = user.Username
			}
		}
	}

	resp := GetRoomResponse{
		ID:           room.ID.String(),
		Title:        room.Title,
		Category:     room.Category,
		CoverURL:     room.CoverURL,
		ChannelName:  room.ChannelName,
		Status:       room.Status,
		StreamerID:   room.StreamerID.String(),
		StreamerName: streamerName,
		StartAt:      formatTime(room.StartAt),
		PeakOnline:   room.PeakOnline,
		TotalViews:   room.TotalViews,
	}

	if room.Status == "live" {
		resp.StreamURL = "rtmp://localhost/live/" + room.ChannelName
		resp.FLVURL = "http://localhost:8080/live/" + room.ChannelName + ".flv"
		resp.HLSURL = "http://localhost:8080/live/" + room.ChannelName + ".m3u8"
	}

	response.Success(c, resp)
}

type RoomListItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	CoverURL    string `json:"cover_url"`
	ChannelName string `json:"channel_name"`
	Status      string `json:"status"`
	StreamerID  string `json:"streamer_id"`
	StartAt     string `json:"start_at"`
	PeakOnline  int    `json:"peak_online"`
	TotalViews  int    `json:"total_views"`
	FLVURL      string `json:"flv_url,omitempty"`
	HLSURL      string `json:"hls_url,omitempty"`
}

func (h *LiveHandler) ListRooms(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	search := c.Query("search")

	query := repository.DB.Model(&models.LiveRoom{}).Preload("Streamer")

	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", "live")
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	var rooms []models.LiveRoom
	query.Order("start_at DESC").Find(&rooms)

	result := make([]RoomListItem, 0)
	for _, room := range rooms {
		streamerName := ""
		if room.StreamerID != uuid.Nil {
			var user models.User
			if err := repository.DB.Select("nickname").Where("id = ?", room.StreamerID).First(&user).Error; err == nil {
				streamerName = user.Nickname
				if streamerName == "" {
					streamerName = user.Username
				}
			}
		}
		item := RoomListItem{
			ID:          room.ID.String(),
			Title:       room.Title,
			Category:    room.Category,
			CoverURL:    room.CoverURL,
			ChannelName: room.ChannelName,
			Status:      room.Status,
			StreamerID:  room.StreamerID.String(),
			StartAt:     formatTime(room.StartAt),
			PeakOnline:  room.PeakOnline,
			TotalViews:  room.TotalViews,
		}
		if room.Status == "live" {
			item.FLVURL = "http://localhost:8080/live/" + room.ChannelName + ".flv"
			item.HLSURL = "http://localhost:8080/live/" + room.ChannelName + ".m3u8"
		}
		result = append(result, item)
	}

	if status == "" || status == "live" {
		var relays []models.RelayStream
		relayQuery := repository.DB.Model(&models.RelayStream{}).Where("status = ?", "running")

		if category != "" {
			relayQuery = relayQuery.Where("category = ?", category)
		}

		if search != "" {
			relayQuery = relayQuery.Where("name ILIKE ?", "%"+search+"%")
		}

		relayQuery.Order("created_at DESC").Find(&relays)

		for _, relay := range relays {
			relayCover := relay.CoverURL
			if relayCover == "" {
				relayCover = "https://placeholder.com/relay-" + relay.Category + ".png"
			}
			startAt := relay.CreatedAt
			if startAt.IsZero() {
				startAt = time.Now()
			}
			result = append(result, RoomListItem{
				ID:          relay.ID.String(),
				Title:       relay.Name,
				Category:    relay.Category,
				CoverURL:    relayCover,
				ChannelName: relay.ChannelName,
				Status:      relay.Status,
				StreamerID:  "relay-" + relay.ID.String()[:8],
				StartAt:     formatTimeFromTime(startAt),
				PeakOnline:  int(relay.PeakOnline),
				TotalViews:  int(relay.ViewCount),
				FLVURL:      "http://localhost:8080/live/" + relay.ChannelName + ".flv",
				HLSURL:      "http://localhost:8080/live/" + relay.ChannelName + ".m3u8",
			})
		}
	}

	response.Success(c, result)
}

type StreamerHandler struct{}

func NewStreamerHandler() *StreamerHandler {
	return &StreamerHandler{}
}

type ApplyStreamerRequest struct {
	Phone string `json:"phone"`
}

func (h *StreamerHandler) Apply(c *gin.Context) {
	userID := c.GetString("user_id")

	var existing models.Streamer
	if err := repository.DB.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		response.BadRequest(c, "you are already a streamer")
		return
	}

	streamKey := generateStreamKey()

	streamer := models.Streamer{
		UserID:     uuid.MustParse(userID),
		StreamKey:  streamKey,
		Status:     "offline",
		IsVerified: false,
	}

	if err := repository.DB.Create(&streamer).Error; err != nil {
		response.Fail(c, "failed to apply")
		return
	}

	response.Success(c, gin.H{
		"message":    "apply successful",
		"stream_key": streamKey,
		"stream_url": "rtmp://localhost/live",
	})
}

func (h *StreamerHandler) GetInfo(c *gin.Context) {
	userID := c.GetString("user_id")

	var streamer models.Streamer
	if err := repository.DB.Where("user_id = ?", userID).First(&streamer).Error; err != nil {
		response.BadRequest(c, "you are not a streamer")
		return
	}

	response.Success(c, gin.H{
		"stream_key":          streamer.StreamKey,
		"status":              streamer.Status,
		"is_verified":         streamer.IsVerified,
		"total_revenue":       streamer.TotalRevenue,
		"follower_count":      streamer.FollowerCount,
		"total_live_duration": streamer.TotalLiveDuration,
	})
}

func (h *StreamerHandler) RefreshStreamKey(c *gin.Context) {
	userID := c.GetString("user_id")

	var streamer models.Streamer
	if err := repository.DB.Where("user_id = ?", userID).First(&streamer).Error; err != nil {
		response.BadRequest(c, "you are not a streamer")
		return
	}

	streamer.StreamKey = generateStreamKey()
	repository.DB.Save(&streamer)

	response.Success(c, gin.H{
		"stream_key": streamer.StreamKey,
	})
}

func generateStreamKey() string {
	return uuid.New().String()[:8] + "-" + uuid.New().String()[:8] + "-" + uuid.New().String()[:8]
}

func ptrTimeNow() *time.Time {
	now := time.Now()
	return &now
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatTimeFromTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
