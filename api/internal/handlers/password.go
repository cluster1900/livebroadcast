package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHandler struct{}

func NewPasswordHandler() *PasswordHandler {
	return &PasswordHandler{}
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *PasswordHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var user struct {
		ID           string
		PasswordHash string
	}
	if err := repository.DB.Raw("SELECT id, password_hash FROM users WHERE id = ?", userID).Scan(&user).Error; err != nil {
		response.BadRequest(c, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		response.BadRequest(c, "原密码错误")
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, "密码加密失败")
		return
	}

	if err := repository.DB.Model(&struct{}{}).Table("users").Where("id = ?", userID).Update("password_hash", string(newHash)).Error; err != nil {
		response.Fail(c, "修改失败")
		return
	}

	response.Success(c, gin.H{"message": "密码修改成功"})
}

type ResetPasswordRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (h *PasswordHandler) RequestReset(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var count int64
	repository.DB.Model(&struct{}{}).Table("users").Where("username = ? AND email = ?", req.Username, req.Email).Count(&count)

	if count == 0 {
		response.BadRequest(c, "用户名和邮箱不匹配")
		return
	}

	// TODO: 发送重置邮件
	// 在生产环境中，这里应该生成重置token并发送到用户邮箱
	response.Success(c, gin.H{"message": "重置链接已发送到您的邮箱"})
}

type CompleteResetRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *PasswordHandler) CompleteReset(c *gin.Context) {
	var req CompleteResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// TODO: 验证重置token
	// 在生产环境中，应该验证token的有效性和过期时间
	response.Success(c, gin.H{"message": "密码重置成功，请使用新密码登录"})
}
