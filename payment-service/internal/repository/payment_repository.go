package repository

import (
	"payment-service/internal/model"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.PaymentSession) error
	FindByID(id string) (*model.PaymentSession, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *model.PaymentSession) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id string) (*model.PaymentSession, error) {
	var payment model.PaymentSession
	err := r.db.Where("id = ?", id).First(&payment).Error
	return &payment, err
}