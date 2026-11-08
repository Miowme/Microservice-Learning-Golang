package service

import (
	"errors"
	"payment-service/internal/model"
	"payment-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(userID, merchantID string, amount float64, description string) (*model.Payment, error)
	GetPaymentByStatus(paymentID string) (*model.Payment, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentService{paymentRepo: paymentRepo}
}

func (s *paymentService) CreatePayment(userID, merchantID string, amount float64, description string) (*model.Payment, error) {
	paymentID := uuid.New().String()
	
	payment := &model.Payment{
		PaymentID: paymentID,
		UserID:    userID,
		Amount:    amount,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.paymentRepo.Create(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *paymentService) GetPaymentByStatus(paymentID string) (*model.Payment, error) {
	payment, err := s.paymentRepo.FindByPaymentID(paymentID)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	return payment, nil
}