package service

import (
	"core-service/internal/model"
	"core-service/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService interface {
	TopUp(userID string, amount float64) (*model.Transaction, float64, error)
	Transfer(fromUserID, toUserID string, amount float64) (*model.Transaction, error)
}

type transactionService struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	db              *gorm.DB
}

func NewTransactionService(walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository, db *gorm.DB) TransactionService {
	return &transactionService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		db:              db,
	}
}

func (s *transactionService) TopUp(userID string, amount float64) (*model.Transaction, float64, error) {
	if amount <= 0 {
		return nil, 0, errors.New("amount must be greater than zero")
	}

	// Menambahkan 'nil' sebagai argumen pertama karena di luar transaksi database
	wallet, err := s.walletRepo.FindByUserID(nil, userID)
	if err != nil {
		return nil, 0, errors.New("wallet not found")
	}
	currentBalance := wallet.Balance

	var transaction *model.Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		err := s.walletRepo.UpdateBalance(tx, userID, amount)
		if err != nil {
			return err
		}

		parsedUserUUID, _ := uuid.Parse(userID)
		desc := "Wallet top-up"
		transaction = &model.Transaction{
			Type:        "topup",
			Amount:      amount,
			ToUserID:    &parsedUserUUID,
			Status:      "completed",
			Description: &desc,
		}

		return s.transactionRepo.Create(tx, transaction)
	})

	if err != nil {
		return nil, 0, err
	}

	newBalance := currentBalance + amount
	return transaction, newBalance, nil
}

func (s *transactionService) Transfer(fromUserID, toUserID string, amount float64) (*model.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if fromUserID == toUserID {
		return nil, errors.New("cannot transfer to the same account")
	}

	// Menambahkan 'nil' sebagai argumen pertama
	senderWallet, err := s.walletRepo.FindByUserID(nil, fromUserID)
	if err != nil || senderWallet.Balance < amount {
		return nil, errors.New("insufficient funds")
	}

	// Menambahkan 'nil' sebagai argumen pertama
	_, err = s.walletRepo.FindByUserID(nil, toUserID)
	if err != nil {
		return nil, errors.New("recipient wallet not found")
	}

	var transaction *model.Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		err = s.walletRepo.UpdateBalance(tx, fromUserID, -amount)
		if err != nil {
			return err
		}

		err = s.walletRepo.UpdateBalance(tx, toUserID, amount)
		if err != nil {
			return err
		}

		parsedFrom, _ := uuid.Parse(fromUserID)
		parsedTo, _ := uuid.Parse(toUserID)
		desc := "Fund transfer"

		transaction = &model.Transaction{
			Type:        "transfer_out",
			Amount:      amount,
			FromUserID:  &parsedFrom,
			ToUserID:    &parsedTo,
			Status:      "completed",
			Description: &desc,
		}

		return s.transactionRepo.Create(tx, transaction)
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}