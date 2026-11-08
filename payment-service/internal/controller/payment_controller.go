package controller

import (
	"net/http"
	"payment-service/internal/model" 
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm" 
)

type PaymentController struct {
	paymentService service.PaymentService
	db             *gorm.DB 
}

func NewPaymentController(paymentService service.PaymentService, db *gorm.DB) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
		db:             db,
	}
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
		"paymentID": payment.PaymentID, 
		"status":    payment.Status,
		"createdAt": payment.CreatedAt, 
	})
}

func (c *PaymentController) GetPaymentStatus(ctx *gin.Context) {
	paymentID := ctx.Param("id")
	var payment model.Payment

	if err := c.db.First(&payment, "payment_id = ?", paymentID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"payment_id": payment.PaymentID,
		"user_id":    payment.UserID,
		"amount":     payment.Amount,
		"status":     payment.Status, 
		"created_at": payment.CreatedAt,
		"updated_at": payment.UpdatedAt,
	})
}