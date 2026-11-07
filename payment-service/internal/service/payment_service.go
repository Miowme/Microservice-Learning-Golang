package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"payment-service/internal/model"
	"payment-service/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(userID, merchantID string, amount float64, description string) (*model.PaymentSession, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	db          *gorm.DB
}

func NewPaymentService(paymentRepo repository.PaymentRepository, db *gorm.DB) PaymentService {
	return &paymentService{paymentRepo: paymentRepo, db: db}
}

func (s *paymentService) CreatePayment(userID, merchantID string, amount float64, description string) (*model.PaymentSession, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	coreServiceURL := os.Getenv("CORE_SERVICE_URL")
	payload := map[string]interface{}{
		"userID": userID,
		"amount": amount,
	}
	jsonBody, _ := json.Marshal(payload)

	resp, err := http.Post(coreServiceURL+"/internal/wallets/deduct", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("gagal memotong saldo di Core Service (saldo tidak cukup atau core service offline)")
	}
	defer resp.Body.Close()

	payment := &model.PaymentSession{
		UserID:      parsedUserUUID,
		MerchantID:  merchantID,
		Amount:      amount,
		Status:      "completed",
		Description: description,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func StartPaymentExpiryWorker(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			now := time.Now()
			result := db.Model(&model.PaymentSession{}).
				Where("status = ? AND expires_at < ?", "pending", now).
				Update("status", "expired")

			if result.Error != nil {
				log.Printf("[Payment Worker Error] Gagal update status expired: %v", result.Error)
			} else if result.RowsAffected > 0 {
				log.Printf("[Payment Worker] Berhasil mengubah %d sesi pembayaran menjadi expired.", result.RowsAffected)
			}
		}
	}()
}