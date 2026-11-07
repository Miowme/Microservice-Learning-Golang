package main

import (
	"log"
	"os"
	"payment-service/internal/config"
	"payment-service/internal/controller"
	"payment-service/internal/model"
	"payment-service/internal/repository"
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := config.ConnectDB()

	db.AutoMigrate(&model.PaymentSession{})

	service.StartPaymentExpiryWorker(db)

	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, db)
	paymentController := controller.NewPaymentController(paymentService)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/transactions/pay", paymentController.Pay)
	}

	port := os.Getenv("PAYMENT_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Payment Service berjalan pada port %s", port)
	r.Run(":" + port)
}