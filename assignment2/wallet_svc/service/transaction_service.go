package service

import (
	"context"
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
		if err == nil {
			//update wallet
			getWallet, err := s.walletRepo.GetWalletByUserID(ctx, userId)
			if err == nil {
				balance := getWallet.Balance - dataTransaction.Amount

				// dataWallet := entity.Wallet{

				// }
			}

		}

	} else {

	}
	// Memanggil CreateWallet dari repository untuk membuat wallet baru
	// createdTransaction, err := s.transactionRepo.CreateTransaction(ctx, transaction)
	// if err != nil {
	// 	return entity.Transaction{}, fmt.Errorf("gagal membuat transaction: %v", err)
	// }
	return entity.Transaction{}, nil
}
