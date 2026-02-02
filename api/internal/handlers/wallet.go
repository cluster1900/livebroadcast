package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/internal/models"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/pkg/response"
)

type WalletHandler struct{}

func NewWalletHandler() *WalletHandler {
	return &WalletHandler{}
}

type RechargeRequest struct {
	Amount int    `json:"amount" binding:"required,min=1"`
	Method string `json:"method" binding:"required"`
}

type RechargeResponse struct {
	TransactionID int64  `json:"transaction_id"`
	Amount        int    `json:"amount"`
	BalanceBefore int    `json:"balance_before"`
	BalanceAfter  int    `json:"balance_after"`
	Method        string `json:"method"`
}

func (h *WalletHandler) RechargeCoins(c *gin.Context) {
	userID := c.GetString("user_id")

	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if req.Amount < 1 || req.Amount > 10000 {
		response.BadRequest(c, "amount must be between 1 and 10000")
		return
	}

	var user models.User
	if err := repository.DB.First(&user, "id = ?", userID).Error; err != nil {
		response.BadRequest(c, "user not found")
		return
	}

	balanceBefore := user.CoinBalance
	balanceAfter := balanceBefore + req.Amount

	user.CoinBalance = balanceAfter
	if err := repository.DB.Save(&user).Error; err != nil {
		response.Fail(c, "failed to recharge")
		return
	}

	coinTx := models.CoinTransaction{
		UserID:       user.ID,
		Amount:       req.Amount,
		BalanceAfter: balanceAfter,
		Type:         "recharge",
		Description:  "Coin recharge via " + req.Method,
	}
	repository.DB.Create(&coinTx)

	response.Success(c, RechargeResponse{
		TransactionID: coinTx.ID,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Method:        req.Method,
	})
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	if err := repository.DB.Select("coin_balance", "level", "exp").First(&user, "id = ?", userID).Error; err != nil {
		response.BadRequest(c, "user not found")
		return
	}

	var levelConfig models.LevelConfig
	repository.DB.First(&levelConfig, "level = ?", user.Level)

	response.Success(c, gin.H{
		"coin_balance": user.CoinBalance,
		"level":        user.Level,
		"exp":          user.Exp,
		"level_name":   levelConfig.LevelName,
		"next_exp":     levelConfig.ExpRequired,
	})
}

func (h *WalletHandler) GetTransactionHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	limit := 20
	offset := 0

	var transactions []models.CoinTransaction
	repository.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions)

	response.Success(c, transactions)
}
