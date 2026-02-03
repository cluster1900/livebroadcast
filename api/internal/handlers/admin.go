package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var stats struct {
		TotalUsers     int   `json:"total_users"`
		TotalStreamers int   `json:"total_streamers"`
		TotalRooms     int   `json:"total_rooms"`
		LiveRooms      int   `json:"live_rooms"`
		TotalRevenue   int64 `json:"total_revenue"`
		PendingReports int   `json:"pending_reports"`
		NewUsersToday  int   `json:"new_users_today"`
	}

	today := time.Now().Format("2006-01-02")

	repository.DB.Raw("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	repository.DB.Raw("SELECT COUNT(*) FROM streamers").Scan(&stats.TotalStreamers)
	repository.DB.Raw("SELECT COUNT(*) FROM live_rooms").Scan(&stats.TotalRooms)
	repository.DB.Raw("SELECT COUNT(*) FROM live_rooms WHERE status = 'live'").Scan(&stats.LiveRooms)
	repository.DB.Raw("SELECT COALESCE(SUM(coin_amount), 0) FROM gift_transactions").Scan(&stats.TotalRevenue)
	repository.DB.Raw("SELECT COUNT(*) FROM user_reports WHERE status = 'pending'").Scan(&stats.PendingReports)
	repository.DB.Raw("SELECT COUNT(*) FROM users WHERE created_at >= ?", today).Scan(&stats.NewUsersToday)

	response.Success(c, stats)
}

func (h *AdminHandler) GetUserList(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var users []struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		Nickname    string `json:"nickname"`
		Level       int    `json:"level"`
		CoinBalance int    `json:"coin_balance"`
		Status      string `json:"status"`
		CreatedAt   string `json:"created_at"`
	}

	repository.DB.Raw(`
		SELECT id, username, nickname, level, coin_balance, status, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT 50
	`).Scan(&users)

	response.Success(c, users)
}

func (h *AdminHandler) BanUser(c *gin.Context) {
	adminID := c.GetString("user_id")
	role := c.GetString("user_role")

	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	userID := c.Param("id")
	if userID == adminID {
		response.BadRequest(c, "不能封禁自己")
		return
	}

	if err := repository.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", "banned").Error; err != nil {
		response.Fail(c, "操作失败")
		return
	}

	response.Success(c, gin.H{"message": "用户已封禁"})
}

func (h *AdminHandler) UnbanUser(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	userID := c.Param("id")

	if err := repository.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", "active").Error; err != nil {
		response.Fail(c, "操作失败")
		return
	}

	response.Success(c, gin.H{"message": "已解封用户"})
}

func (h *AdminHandler) GetRoomList(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var rooms []struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		StreamerID string `json:"streamer_id"`
		Streamer   string `json:"streamer"`
		Status     string `json:"status"`
		PeakOnline int    `json:"peak_online"`
		TotalViews int    `json:"total_views"`
		CreatedAt  string `json:"created_at"`
	}

	repository.DB.Raw(`
		SELECT r.id, r.title, r.streamer_id, u.nickname as streamer, 
		       r.status, r.peak_online, r.total_views, r.created_at
		FROM live_rooms r
		JOIN users u ON u.id = r.streamer_id
		ORDER BY r.created_at DESC
		LIMIT 50
	`).Scan(&rooms)

	response.Success(c, rooms)
}

func (h *AdminHandler) BanRoom(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	roomID := c.Param("id")
	reason := c.Query("reason")

	if err := repository.DB.Model(&models.LiveRoom{}).Where("id = ?", roomID).Updates(map[string]interface{}{
		"status": "banned",
	}).Error; err != nil {
		response.Fail(c, "操作失败")
		return
	}

	response.Success(c, gin.H{"message": "直播间已封禁", "reason": reason})
}

func (h *AdminHandler) GetGiftList(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var gifts []models.Gift
	repository.DB.Order("sort_order").Find(&gifts)

	response.Success(c, gifts)
}

func (h *AdminHandler) CreateGift(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var gift models.Gift
	if err := c.ShouldBindJSON(&gift); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := repository.DB.Create(&gift).Error; err != nil {
		response.Fail(c, "创建失败")
		return
	}

	response.Success(c, gin.H{"message": "创建成功", "data": gift})
}

func (h *AdminHandler) UpdateGift(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	giftID := c.Param("id")

	var gift models.Gift
	if err := repository.DB.First(&gift, "id = ?", giftID).Error; err != nil {
		response.Fail(c, "礼物不存在")
		return
	}

	if err := c.ShouldBindJSON(&gift); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	repository.DB.Save(&gift)

	response.Success(c, gin.H{"message": "更新成功"})
}

func (h *AdminHandler) DeleteGift(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	giftID := c.Param("id")

	if err := repository.DB.Delete(&models.Gift{}, "id = ?", giftID).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

func (h *AdminHandler) GetSensitiveWords(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var words []models.SensitiveWord
	repository.DB.Find(&words)

	response.Success(c, words)
}

func (h *AdminHandler) AddSensitiveWord(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var word models.SensitiveWord
	if err := c.ShouldBindJSON(&word); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := repository.DB.Create(&word).Error; err != nil {
		response.Fail(c, "添加失败")
		return
	}

	response.Success(c, gin.H{"message": "添加成功"})
}

func (h *AdminHandler) DeleteSensitiveWord(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	wordID := c.Param("id")

	if err := repository.DB.Delete(&models.SensitiveWord{}, "id = ?", wordID).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

func (h *AdminHandler) GetSystemConfig(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var configs []models.SystemConfig
	repository.DB.Find(&configs)

	response.Success(c, configs)
}

func (h *AdminHandler) UpdateSystemConfig(c *gin.Context) {
	role := c.GetString("user_role")
	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := repository.DB.Model(&models.SystemConfig{}).Where("key = ?", req.Key).Update("value", req.Value).Error; err != nil {
		response.Fail(c, "更新失败")
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}
