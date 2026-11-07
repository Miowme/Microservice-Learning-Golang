package controller

import (
	"net/http"
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService service.PaymentService
}

func NewPaymentController(paymentService service.PaymentService) *PaymentController {
	return &PaymentController{paymentService: paymentService}
}

func (c *PaymentController) Pay(ctx *gin.Context) {
	var req struct {
		UserID      string  `json:"userID" binding:"required"`
		MerchantID  string  `json:"merchantID" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		Description string  `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := c.paymentService.CreatePayment(req.UserID, req.MerchantID, req.Amount, req.Description)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Payment successful",
		"paymentID": payment.ID,
		"status":    payment.Status,
		"expiresAt": payment.ExpiresAt,
	})
}