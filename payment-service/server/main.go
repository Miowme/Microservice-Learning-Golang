package main

import (
	"log"
	"os"

	"payment-service/internal/config"
	"payment-service/internal/controller" 
	"payment-service/internal/model"
	"payment-service/internal/rabbitmq"
	"payment-service/internal/repository"
	"payment-service/internal/service" 

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := config.ConnectDB()

	db.AutoMigrate(&model.Payment{})

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	go func() {
		log.Println("[Payment Service] Memulai RabbitMQ Consumer worker...")
		rabbitmq.StartConsumer(rabbitURL, db)
	}()
	
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo) 
	paymentController := controller.NewPaymentController(paymentService, db)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"status": "Payment Service is running"})
		})

		v1.GET("/payments/:id", paymentController.GetPaymentStatus)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Payment Service berjalan pada port %s", port)
	r.Run(":" + port)
}