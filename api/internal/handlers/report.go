package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type ReportHandler struct{}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

type CreateReportRequest struct {
	ReportedID string `json:"reported_id" binding:"required"`
	RoomID     string `json:"room_id"`
	Type       string `json:"type" binding:"required"`
	Reason     string `json:"reason" binding:"required,max=500"`
}

func (h *ReportHandler) CreateReport(c *gin.Context) {
	reporterID := c.GetString("user_id")
	if reporterID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	reporterUUID := uuid.MustParse(reporterID)
	reportedUUID, err := uuid.Parse(req.ReportedID)
	if err != nil {
		response.BadRequest(c, "被举报者ID无效")
		return
	}

	if reporterID == req.ReportedID {
		response.BadRequest(c, "不能举报自己")
		return
	}

	var roomUUID *uuid.UUID
	if req.RoomID != "" {
		ru, err := uuid.Parse(req.RoomID)
		if err != nil {
			response.BadRequest(c, "直播间ID无效")
			return
		}
		roomUUID = &ru
	}

	report := models.UserReport{
		ID:         uuid.New(),
		ReporterID: reporterUUID,
		ReportedID: reportedUUID,
		RoomID:     roomUUID,
		Type:       req.Type,
		Reason:     req.Reason,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	if err := repository.DB.Create(&report).Error; err != nil {
		response.Fail(c, "举报失败")
		return
	}

	response.Success(c, gin.H{"message": "举报已提交，感谢您的反馈"})
}

func (h *ReportHandler) GetMyReports(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	userUUID := uuid.MustParse(userID)
	var reports []models.UserReport
	repository.DB.Where("reporter_id = ?", userUUID).Order("created_at DESC").Find(&reports)

	response.Success(c, reports)
}

func (h *ReportHandler) GetPendingReports(c *gin.Context) {
	role := c.GetString("user_role")

	if role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	var reports []struct {
		ID         string `json:"id"`
		ReporterID string `json:"reporter_id"`
		Reporter   string `json:"reporter"`
		ReportedID string `json:"reported_id"`
		Reported   string `json:"reported"`
		RoomID     string `json:"room_id"`
		Type       string `json:"type"`
		Reason     string `json:"reason"`
		Status     string `json:"status"`
		CreatedAt  string `json:"created_at"`
	}

	repository.DB.Raw(`
		SELECT 
			r.id, r.reporter_id, u1.nickname as reporter,
			r.reported_id, u2.nickname as reported,
			r.room_id, r.type, r.reason, r.status, r.created_at
		FROM user_reports r
		JOIN users u1 ON u1.id = r.reporter_id
		JOIN users u2 ON u2.id = r.reported_id
		WHERE r.status = 'pending'
		ORDER BY r.created_at DESC
	`).Scan(&reports)

	response.Success(c, reports)
}

type HandleReportRequest struct {
	Status     string `json:"status" binding:"required"`
	HandleNote string `json:"handle_note"`
}

func (h *ReportHandler) HandleReport(c *gin.Context) {
	adminID := c.GetString("user_id")
	role := c.GetString("user_role")

	if adminID == "" || role != "admin" {
		response.Forbidden(c, "需要管理员权限")
		return
	}

	reportID := c.Param("id")
	var req HandleReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	adminUUID := uuid.MustParse(adminID)
	now := time.Now()

	if err := repository.DB.Model(&models.UserReport{}).Where("id = ?", reportID).Updates(map[string]interface{}{
		"status":      req.Status,
		"handle_by":   adminUUID,
		"handle_note": req.HandleNote,
		"handle_at":   &now,
	}).Error; err != nil {
		response.Fail(c, "处理失败")
		return
	}

	response.Success(c, gin.H{"message": "处理成功"})
}
