package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type LikeHandler struct{}

func NewLikeHandler() *LikeHandler {
	return &LikeHandler{}
}

func (h *LikeHandler) LikeRoom(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	roomID := c.Param("room_id")
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		response.BadRequest(c, "直播间ID无效")
		return
	}

	userUUID := uuid.MustParse(userID)

	var like models.RoomLike
	if err := repository.DB.Where("user_id = ? AND room_id = ?", userUUID, roomUUID).First(&like).Error; err == nil {
		response.Fail(c, "您已经点赞过了")
		return
	}

	like = models.RoomLike{
		ID:     uuid.New(),
		UserID: userUUID,
		RoomID: roomUUID,
	}

	if err := repository.DB.Create(&like).Error; err != nil {
		response.Fail(c, "点赞失败")
		return
	}

	response.Success(c, gin.H{"message": "点赞成功"})
}

func (h *LikeHandler) UnlikeRoom(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	roomID := c.Param("room_id")
	roomUUID := uuid.MustParse(roomID)
	userUUID := uuid.MustParse(userID)

	if err := repository.DB.Where("user_id = ? AND room_id = ?", userUUID, roomUUID).Delete(&models.RoomLike{}).Error; err != nil {
		response.Fail(c, "取消点赞失败")
		return
	}

	response.Success(c, gin.H{"message": "取消点赞成功"})
}

func (h *LikeHandler) GetLikeCount(c *gin.Context) {
	roomID := c.Param("room_id")
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		response.BadRequest(c, "直播间ID无效")
		return
	}

	var count int64
	repository.DB.Model(&models.RoomLike{}).Where("room_id = ?", roomUUID).Count(&count)

	response.Success(c, gin.H{"count": count})
}

func (h *LikeHandler) HasLiked(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	roomID := c.Param("room_id")
	roomUUID := uuid.MustParse(roomID)
	userUUID := uuid.MustParse(userID)

	var count int64
	repository.DB.Model(&models.RoomLike{}).Where("user_id = ? AND room_id = ?", userUUID, roomUUID).Count(&count)

	response.Success(c, gin.H{"has_liked": count > 0})
}
