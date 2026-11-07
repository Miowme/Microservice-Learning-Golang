package main

import (
	"log"
	"os"

	"core-service/internal/config"
	"core-service/internal/controller"
	"core-service/internal/model"
	"core-service/internal/repository"
	"core-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := config.ConnectDB()

	db.AutoMigrate(&model.User{}, &model.Wallet{}, &model.Transaction{})

	//Repository
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	//Service (pastikan parameter sesuai: repo-repo lalu db)
	authService := service.NewAuthService(userRepo, walletRepo, db)
	walletService := service.NewWalletService(walletRepo, db)
	transactionService := service.NewTransactionService(walletRepo, transactionRepo, db)

	//Controller
	authController := controller.NewAuthController(authService)
	walletController := controller.NewWalletController(walletService)
	transactionController := controller.NewTransactionController(transactionService)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", authController.Register)
		v1.POST("/login", authController.Login)

		v1.POST("/wallets/topup", transactionController.TopUp)
		v1.POST("/wallets/transfer", transactionController.Transfer)

		v1.POST("/internal/wallets/deduct", walletController.DeductBalanceInternal)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Core Service berjalan pada port %s", port)
	r.Run(":" + port)
}