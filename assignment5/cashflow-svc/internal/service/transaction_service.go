package service

import (
	"assignment5/cashflow-svc/internal/entity"
	"assignment5/cashflow-svc/internal/model"
	"assignment5/cashflow-svc/internal/repository/postgres_gorm"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type ITransactionService interface {
	TransferWallet(ctx context.Context, request model.TransferWalletRequest) (model.TransferWalletResponse, error)
}

type transactionService struct {
	db              *gorm.DB
	walletRepo      postgres_gorm.IWalletRepository
	categoryRepo    postgres_gorm.ICategoryRepository
	transactionRepo postgres_gorm.ITransactionRepository
}

func NewTransactionService(
	db *gorm.DB,
	walletRepo postgres_gorm.IWalletRepository,
	categoryRepo postgres_gorm.ICategoryRepository,
	transactionRepo postgres_gorm.ITransactionRepository,
) *transactionService {
	return &transactionService{
		db:              db,
		walletRepo:      walletRepo,
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) TransferWallet(ctx context.Context, request model.TransferWalletRequest) (model.TransferWalletResponse, error) {
	tx := s.db.WithContext(ctx).Begin()

	//Get Wallet Sender
	senderWallet, err := s.walletRepo.GetWalletByID(ctx, int(request.FromID))
	if err != nil {
		tx.Rollback()
		return model.TransferWalletResponse{}, errors.New("Sender's wallet id not found")
	} else {
		recipientWallet, err := s.walletRepo.GetWalletByID(ctx, int(request.ToID))
		if err != nil {
			tx.Rollback()
			return model.TransferWalletResponse{}, errors.New("Recipient Wallet ID not found")
		}

		if senderWallet.UserID != recipientWallet.UserID {
			tx.Rollback()
			return model.TransferWalletResponse{}, errors.New("Different Wallet Users")
		}

		if senderWallet.Balance < request.Nominal {
			tx.Rollback()
			return model.TransferWalletResponse{}, errors.New("Insufficient balance to carry out transactions. Your current balance: " + fmt.Sprint("%f", senderWallet.Balance))
		}

		//create transaction out
		var transactionOut = &entity.Transaction{
			WalletID:        int(request.FromID),
			Nominal:         request.Nominal,
			Type:            "Withdrawal",
			FromWalletID:    int(request.ToID),
			ToWalletID:      senderWallet.ID,
			TransactionDate: request.TransactionDate,
			CategoryID:      nil,
		}
		saveTransactionOut, err := s.transactionRepo.CreateTransaction(ctx, transactionOut)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return model.TransferWalletResponse{}, err
		}

		//update balance wallet sender
		var lastBalanceSender float64
		lastBalanceSender = senderWallet.Balance - request.Nominal
		updateWalletSender := entity.Wallet{
			ID:      senderWallet.ID,
			Balance: lastBalanceSender,
			Name:    senderWallet.Name,
		}
		_, err = s.walletRepo.UpdateWallet(ctx, senderWallet.ID, updateWalletSender)
		if err != nil {
			tx.Rollback()
			return model.TransferWalletResponse{}, err
		}

		//Create transaction IN
		transactionIn := &entity.Transaction{
			WalletID:        int(request.ToID),
			Nominal:         request.Nominal,
			Type:            "Deposit",
			FromWalletID:    int(request.FromID),
			ToWalletID:      int(request.ToID),
			TransactionDate: request.TransactionDate,
		}
		saveTransactionIn, err := s.transactionRepo.CreateTransaction(ctx, transactionIn)
		if err != nil {
			tx.Rollback()
			return model.TransferWalletResponse{}, err
		}

		//update balance wallet receiver
		var lastBalanceReceiver float64
		lastBalanceReceiver = recipientWallet.Balance + request.Nominal
		updateWalletReceiver := entity.Wallet{
			ID:      recipientWallet.ID,
			Balance: lastBalanceReceiver,
			Name:    recipientWallet.Name,
		}
		_, err = s.walletRepo.UpdateWallet(ctx, recipientWallet.ID, updateWalletReceiver)
		if err != nil {
			tx.Rollback()
			return model.TransferWalletResponse{}, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return model.TransferWalletResponse{}, err
		}

		return model.TransferWalletResponse{
			SenderWalletID:   senderWallet.ID,
			ReceiverWalletID: recipientWallet.ID,
			Nominal:          request.Nominal,
			TransactionDate:  request.TransactionDate,
			TransactionIn:    saveTransactionIn,
			TransactionOut:   saveTransactionOut,
		}, err
	}
}
