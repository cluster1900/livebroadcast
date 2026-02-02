package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/jwt"
	"github.com/huya_live/api/pkg/redis"
	"github.com/huya_live/api/pkg/response"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthHandler struct {
	jwtManager *jwt.Manager
}

func NewAuthHandler(jwtManager *jwt.Manager) *AuthHandler {
	return &AuthHandler{jwtManager: jwtManager}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	var existingUser models.User
	if err := repository.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		response.BadRequest(c, "username already exists")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, "failed to hash password")
		return
	}

	user := models.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Nickname:     req.Nickname,
		Phone:        req.Phone,
		Email:        req.Email,
		Level:        1,
		Exp:          0,
		CoinBalance:  0,
		Status:       "active",
	}

	if err := repository.DB.Create(&user).Error; err != nil {
		response.Fail(c, "failed to create user")
		return
	}

	accessToken, _ := h.jwtManager.GenerateAccessToken(user.ID.String(), user.Username, "user", user.Level)
	refreshToken, _ := h.jwtManager.GenerateRefreshToken(user.ID.String())

	redis.Set(c.Request.Context(), "refresh:"+user.ID.String(), refreshToken, time.Duration(7*24*time.Hour))

	response.Success(c, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	var user models.User
	if err := repository.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response.BadRequest(c, "invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.BadRequest(c, "invalid username or password")
		return
	}

	if user.Status != "active" {
		response.BadRequest(c, "account is disabled")
		return
	}

	now := time.Now()
	repository.DB.Model(&user).Update("last_login_at", &now)

	accessToken, _ := h.jwtManager.GenerateAccessToken(user.ID.String(), user.Username, "user", user.Level)
	refreshToken, _ := h.jwtManager.GenerateRefreshToken(user.ID.String())

	redis.Set(c.Request.Context(), "refresh:"+user.ID.String(), refreshToken, time.Duration(7*24*time.Hour))

	response.Success(c, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid refresh token")
		return
	}

	storedToken, err := redis.Get(c.Request.Context(), "refresh:"+claims.UserID)
	if err != nil || storedToken != req.RefreshToken {
		response.Unauthorized(c, "refresh token mismatch")
		return
	}

	var user models.User
	if err := repository.DB.First(&user, "id = ?", claims.UserID).Error; err != nil {
		response.BadRequest(c, "user not found")
		return
	}

	newAccessToken, _ := h.jwtManager.GenerateAccessToken(user.ID.String(), user.Username, "user", user.Level)

	response.Success(c, gin.H{
		"access_token": newAccessToken,
	})
}

type ProfileResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
	Level       int    `json:"level"`
	Exp         int64  `json:"exp"`
	CoinBalance int    `json:"coin_balance"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var user models.User
	if err := repository.DB.First(&user, "id = ?", userID).Error; err != nil {
		response.BadRequest(c, "user not found")
		return
	}

	response.Success(c, ProfileResponse{
		ID:          user.ID.String(),
		Username:    user.Username,
		Nickname:    user.Nickname,
		AvatarURL:   user.AvatarURL,
		Level:       user.Level,
		Exp:         user.Exp,
		CoinBalance: user.CoinBalance,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}

	if err := repository.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		response.Fail(c, "failed to update profile")
		return
	}

	response.Success(c, gin.H{
		"message": "profile updated successfully",
	})
}
