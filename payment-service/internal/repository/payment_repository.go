package repository

import (
	"payment-service/internal/model"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	FindByPaymentID(paymentID string) (*model.Payment, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByPaymentID(paymentID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.First(&payment, "payment_id = ?", paymentID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}