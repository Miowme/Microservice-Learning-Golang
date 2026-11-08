package model

import (
	"time"
)

type Payment struct {
	PaymentID string    `gorm:"primaryKey" json:"payment_id"`
	UserID    string    `json:"user_id"`
	MerchantID string   `json:"merchant_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Payment) TableName() string {
	return "payment_sessions"
}