package controller

import (
	"core-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	txService service.TransactionService
}

func NewTransactionController(txService service.TransactionService) *TransactionController {
	return &TransactionController{txService: txService}
}

func (c *TransactionController) TopUp(ctx *gin.Context) {
	userID := ctx.GetString("userID") 

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, newBalance, err := c.txService.TopUp(userID, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Top-up successful",
		"newBalance":    newBalance,
		"transactionID": tx.ID,
	})
}

func (c *TransactionController) Transfer(ctx *gin.Context) {
	userID := ctx.GetString("userID") // Pengirim dari JWT Middleware

	var req struct {
		ToUserID string  `json:"toUserID" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := c.txService.Transfer(userID, req.ToUserID, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Transfer successful",
		"transactionID": tx.ID,
		"amount":        tx.Amount,
	})
}