package repository

import (
    "core-service/internal/model"
    "gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaction *model.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *gorm.DB, transaction *model.Transaction) error {
	dbConn := r.db
	if tx != nil {
		dbConn = tx
	}
	return dbConn.Create(transaction).Error
}