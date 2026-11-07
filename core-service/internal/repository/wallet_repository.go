package repository

import (
	"core-service/internal/model"
	"gorm.io/gorm"
)

type WalletRepository interface {
	FindByUserID(tx *gorm.DB, userID string) (*model.Wallet, error)
	FindByUserIDForUpdate(tx *gorm.DB, userID string) (*model.Wallet, error)
	UpdateBalance(tx *gorm.DB, userID string, amount float64) error
	Create(tx *gorm.DB, wallet *model.Wallet) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *walletRepository) FindByUserID(tx *gorm.DB, userID string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.getDB(tx).Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

func (r *walletRepository) FindByUserIDForUpdate(tx *gorm.DB, userID string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.getDB(tx).Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

func (r *walletRepository) UpdateBalance(tx *gorm.DB, userID string, amount float64) error {
	return r.getDB(tx).Model(&model.Wallet{}).
		Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *walletRepository) Create(tx *gorm.DB, wallet *model.Wallet) error {
	return r.getDB(tx).Create(wallet).Error
}