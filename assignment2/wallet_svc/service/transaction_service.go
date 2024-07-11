package service

import (
	"context"
	"fmt"
	"log"
	"wallet_svc/entity"
)

type ITransactionService interface {
	CreateTransaction(ctx context.Context, transaction *entity.TransactionRequest) (entity.Transaction, error)
}

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error)
}

type transactionService struct {
	transactionRepo ITransactionRepository
	walletRepo      IWalletRepository
}

func NewTransactionService(transactionRepo ITransactionRepository, walletRepo IWalletRepository) ITransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, transaction *entity.TransactionRequest) (entity.Transaction, error) {
	if transaction.Type == "topup" {
		//create transaction
		userId := transaction.ToID
		dataTransaction := entity.Transaction{
			UserID:   userId,
			Amount:   transaction.Amount,
			Category: "IN",
			Type:     transaction.Type,
		}
		_, err := s.transactionRepo.CreateTransaction(ctx, &dataTransaction)
		if err != nil {
			return dataTransaction, fmt.Errorf("gagal create transaction: %v", err)
		}
		//update wallet
		getWallet, err := s.walletRepo.GetWalletByUserID(ctx, userId)
		if err != nil {
			return dataTransaction, fmt.Errorf("gagal get wallet: %v", err)
		}
		balance := getWallet.Balance + dataTransaction.Amount
		dataWallet := entity.Wallet{
			ID:      getWallet.ID,
			Balance: balance,
		}
		updatedWallet, err := s.walletRepo.UpdateWallet(ctx, getWallet.ID, dataWallet)
		log.Println(updatedWallet)
		if err != nil {
			return dataTransaction, fmt.Errorf("gagal memperbarui wallet: %v", err)
		}
	} else {
		//untuk mengurangi balance wallet
		dataTransaction := entity.Transaction{
			UserID:   transaction.FromID,
			Amount:   transaction.Amount,
			Category: "OUT",
			Type:     transaction.Type,
		}
		_, err := s.transactionRepo.CreateTransaction(ctx, &dataTransaction)
		if err != nil {
			return dataTransaction, fmt.Errorf("gagal create transaction: %v", err)
		}

		getWalletFrom, err := s.walletRepo.GetWalletByUserID(ctx, transaction.FromID)
		if err != nil {
			return dataTransaction, fmt.Errorf("gagal get wallet: %v", err)
		}
		balance := getWalletFrom.Balance - dataTransaction.Amount
		dataWallet := entity.Wallet{
			ID:      getWalletFrom.ID,
			Balance: balance,
		}
		updatedWallet, err := s.walletRepo.UpdateWallet(ctx, getWalletFrom.ID, dataWallet)
		log.Println(updatedWallet)
		if err != nil {
			return dataTransaction, fmt.Errorf("gagal mengurangi balance wallet: %v", err)
		}

		//untuk menambahakan balance wallet
		dtTransaction := entity.Transaction{
			UserID:   transaction.ToID,
			Amount:   dataTransaction.Amount,
			Category: "IN",
			Type:     transaction.Type,
		}
		createdTrans, err := s.transactionRepo.CreateTransaction(ctx, &dtTransaction)
		log.Println(createdTrans)
		if err != nil {
			return dataTransaction, fmt.Errorf("gagal create transaction: %v", err)
		}

		getWalletTo, err := s.walletRepo.GetWalletByUserID(ctx, transaction.ToID)
		if err != nil {
			return dtTransaction, fmt.Errorf("gagal get wallet: %v", err)
		}
		balanceTo := getWalletTo.Balance + dataTransaction.Amount
		dtWallet := entity.Wallet{
			ID:      getWalletTo.ID,
			Balance: balanceTo,
		}
		updatedWalletTo, err := s.walletRepo.UpdateWallet(ctx, getWalletTo.ID, dtWallet)
		log.Println(updatedWalletTo)
		if err != nil {
			return dataTransaction, fmt.Errorf("gagal mengurangi balance wallet: %v", err)
		}
	}
	return entity.Transaction{}, nil
}
