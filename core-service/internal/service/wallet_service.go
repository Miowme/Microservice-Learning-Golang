package service

import (
	"core-service/internal/model"
	"core-service/internal/repository"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type WalletService interface {
	GetBalance(userID string) (*model.Wallet, error)
	DeductBalance(userID string, amount float64) error
}

type walletService struct {
	walletRepo repository.WalletRepository
	db         *gorm.DB
}

func NewWalletService(walletRepo repository.WalletRepository, db *gorm.DB) WalletService {
	return &walletService{walletRepo: walletRepo, db: db}
}

func (s *walletService) GetBalance(userID string) (*model.Wallet, error) {
	wallet, err := s.walletRepo.FindByUserID(nil, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) DeductBalance(userID string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.walletRepo.FindByUserIDForUpdate(tx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWalletNotFound
			}
			return err
		}

		if wallet.Balance < amount {
			return ErrInsufficientFunds
		}

		return s.walletRepo.UpdateBalance(tx, userID, -amount)
	})
}