package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type SocialHandler struct{}

func NewSocialHandler() *SocialHandler {
	return &SocialHandler{}
}

type FollowRequest struct {
	StreamerID string `json:"streamer_id" binding:"required"`
}

func (h *SocialHandler) Follow(c *gin.Context) {
	userID := c.GetString("user_id")
	var req FollowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if req.StreamerID == userID {
		response.BadRequest(c, "cannot follow yourself")
		return
	}

	var streamer models.Streamer
	if err := repository.DB.Where("user_id = ?", req.StreamerID).First(&streamer).Error; err != nil {
		response.BadRequest(c, "streamer not found")
		return
	}

	var existing models.FanRelation
	if err := repository.DB.Where("user_id = ? AND streamer_id = ?", userID, req.StreamerID).First(&existing).Error; err == nil {
		response.BadRequest(c, "already following this streamer")
		return
	}

	now := time.Now()
	fanRelation := models.FanRelation{
		UserID:        uuid.MustParse(userID),
		StreamerID:    uuid.MustParse(req.StreamerID),
		FanLevel:      1,
		LoyaltyPoints: 0,
		FollowedAt:    now,
	}
	repository.DB.Create(&fanRelation)

	streamer.FollowerCount++
	repository.DB.Save(&streamer)

	response.Success(c, gin.H{
		"message":        "followed successfully",
		"streamer_id":    req.StreamerID,
		"fan_level":      1,
		"follower_count": streamer.FollowerCount,
	})
}

func (h *SocialHandler) Unfollow(c *gin.Context) {
	userID := c.GetString("user_id")
	var req FollowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	var fanRelation models.FanRelation
	if err := repository.DB.Where("user_id = ? AND streamer_id = ?", userID, req.StreamerID).First(&fanRelation).Error; err != nil {
		response.BadRequest(c, "not following this streamer")
		return
	}

	repository.DB.Delete(&fanRelation)

	var streamer models.Streamer
	if err := repository.DB.Where("user_id = ?", req.StreamerID).First(&streamer).Error; err == nil {
		if streamer.FollowerCount > 0 {
			streamer.FollowerCount--
			repository.DB.Save(&streamer)
		}
	}

	response.Success(c, gin.H{
		"message":     "unfollowed successfully",
		"streamer_id": req.StreamerID,
	})
}

func (h *SocialHandler) GetFollowings(c *gin.Context) {
	userID := c.GetString("user_id")

	var relations []models.FanRelation
	repository.DB.Preload("Streamer").Where("user_id = ?", userID).Find(&relations)

	result := make([]gin.H, 0)
	for _, rel := range relations {
		var user models.User
		repository.DB.Select("nickname", "avatar_url").First(&user, "id = ?", rel.StreamerID)

		result = append(result, gin.H{
			"streamer_id":    rel.StreamerID,
			"nickname":       user.Nickname,
			"avatar":         user.AvatarURL,
			"fan_level":      rel.FanLevel,
			"loyalty_points": rel.LoyaltyPoints,
			"followed_at":    rel.FollowedAt,
		})
	}

	response.Success(c, result)
}

func (h *SocialHandler) GetFollowers(c *gin.Context) {
	streamerID := c.Param("streamer_id")

	var relations []models.FanRelation
	repository.DB.Preload("User").Where("streamer_id = ?", streamerID).Find(&relations)

	result := make([]gin.H, 0)
	for _, rel := range relations {
		var user models.User
		repository.DB.Select("nickname", "avatar_url").First(&user, "id = ?", rel.UserID)

		result = append(result, gin.H{
			"user_id":        rel.UserID,
			"nickname":       user.Nickname,
			"avatar":         user.AvatarURL,
			"fan_level":      rel.FanLevel,
			"loyalty_points": rel.LoyaltyPoints,
		})
	}

	response.Success(c, result)
}
