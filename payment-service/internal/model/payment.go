package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentSession struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	MerchantID  string    `gorm:"type:varchar(100);not null"`
	Amount      float64   `gorm:"type:numeric(18,2);not null"`
	Status      string    `gorm:"type:varchar(20);not null;default:'pending'"` // pending, completed, expired
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	ExpiresAt   time.Time `gorm:"not null"`
}

func (PaymentSession) TableName() string {
	return "payment_sessions"
}