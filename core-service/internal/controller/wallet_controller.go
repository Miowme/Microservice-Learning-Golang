package controller

import (
	"core-service/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	walletService service.WalletService
}

func NewWalletController(walletService service.WalletService) *WalletController {
	return &WalletController{walletService: walletService}
}

type deductBalanceRequest struct {
	UserID string  `json:"userID" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (c *WalletController) InternalDeductBalance(ctx *gin.Context) {
	var req deductBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.walletService.DeductBalance(req.UserID, req.Amount)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, gin.H{"message": "balance deducted successfully"})
	case errors.Is(err, service.ErrInvalidAmount):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrWalletNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrInsufficientFunds):
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func (c *WalletController) DeductBalanceInternal(ctx *gin.Context) {
	var req struct {
		UserID string  `json:"userID" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.walletService.DeductBalance(req.UserID, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Balance deducted successfully"})
}