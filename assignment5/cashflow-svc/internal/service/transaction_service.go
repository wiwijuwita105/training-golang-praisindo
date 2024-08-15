package service

import (
	"assignment5/cashflow-svc/internal/config"
	"assignment5/cashflow-svc/internal/entity"
	"assignment5/cashflow-svc/internal/model"
	"assignment5/cashflow-svc/internal/repository/postgres_gorm"
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type ITransactionService interface {
	TransferWallet(ctx context.Context, request model.TransferWalletRequest) (model.TransferWalletResponse, error)
	CreateTransaction(ctx context.Context, request model.TransactionRequest) (entity.Transaction, error)
	GetTransactions(ctx context.Context, request model.FilterTransactionRequest) ([]model.TransactionResponse, error)
	GetLastTransactions(ctx context.Context, request model.LastTransactionRequest) ([]model.TransactionResponse, error)
	GetCashflowReport(ctx context.Context, request model.CashFlowReportRequest) (model.CashFlowReportResponse, error)
	GetSummaryCategory(ctx context.Context, request model.SummaryCategoryRequest) ([]model.SummaryCategoryResponse, error)
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
			return model.TransferWalletResponse{}, errors.New("Insufficient balance to carry out transactions")
		}

		//create transaction out
		var transactionOut = &entity.Transaction{
			WalletID:        int(request.FromID),
			Nominal:         request.Nominal,
			Type:            config.EXPENSE,
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
		var lastBalanceSender = senderWallet.Balance - request.Nominal
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
			Type:            config.INCOME,
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
		var lastBalanceReceiver = recipientWallet.Balance + request.Nominal
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

func (s *transactionService) CreateTransaction(ctx context.Context, request model.TransactionRequest) (entity.Transaction, error) {
	tx := s.db.WithContext(ctx).Begin()

	getCategory, err := s.categoryRepo.GetCategoryByID(ctx, int(request.CategoryID))
	if err != nil {
		tx.Rollback()
		return entity.Transaction{}, err
	}

	getWallet, err := s.walletRepo.GetWalletByID(ctx, int(request.WalletID))
	if err != nil {
		tx.Rollback()
		return entity.Transaction{}, err
	}

	if getWallet.Balance < request.Nominal && getCategory.Type == config.EXPENSE {
		tx.Rollback()
		return entity.Transaction{}, errors.New("Insufficient balance to carry out transactions")
	}

	var idCategory = getCategory.ID
	var inputTransaction = &entity.Transaction{
		WalletID:        int(request.WalletID),
		CategoryID:      &idCategory,
		Type:            getCategory.Type,
		Nominal:         request.Nominal,
		TransactionDate: time.Now(),
	}
	insertTransaction, err := s.transactionRepo.CreateTransaction(ctx, inputTransaction)
	if err != nil {
		tx.Rollback()
		return entity.Transaction{}, err
	}

	var lastBalance float64
	if getCategory.Type == config.INCOME {
		lastBalance = getWallet.Balance + request.Nominal
	} else {
		lastBalance = getWallet.Balance - request.Nominal
	}

	inputUpdateWallet := entity.Wallet{
		Balance: lastBalance,
	}
	_, err = s.walletRepo.UpdateWallet(ctx, getWallet.ID, inputUpdateWallet)
	if err != nil {
		tx.Rollback()
		return entity.Transaction{}, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return entity.Transaction{}, err
	}
	return insertTransaction, nil
}

func (s *transactionService) GetTransactions(ctx context.Context, request model.FilterTransactionRequest) ([]model.TransactionResponse, error) {
	filter := model.FilterTransaction{}
	filter.StartTime = request.StartTime
	filter.EndTime = request.EndTime
	if request.WalletID != 0 {
		filter.WalletID = []int32{request.WalletID}
	} else {
		//get wallet id by userID
		getWallets, err := s.walletRepo.GetWalletByUserID(ctx, int(request.UserID))
		if err != nil {
			return nil, err
		}
		for _, wallet := range getWallets {
			filter.WalletID = append(filter.WalletID, int32(wallet.ID))
		}
	}

	getTransaction, err := s.transactionRepo.GetAllTransactions(ctx, filter)
	if err != nil {
		return nil, err
	}

	var transactions []model.TransactionResponse
	for _, record := range getTransaction {
		var categoryID int32
		var categoryName string
		if record.CategoryID != nil {
			categoryID = int32(*record.CategoryID)
			categoryName = record.Category.Name
		} else {
			categoryID = 0
			categoryName = ""
		}
		transactions = append(transactions, model.TransactionResponse{
			ID:              int32(record.ID),
			TransactionDate: record.TransactionDate,
			Type:            record.Type,
			Nominal:         record.Nominal,
			WalletID:        int32(record.WalletID),
			WalletName:      record.Wallet.Name,
			CategoryID:      categoryID,
			CategoryName:    categoryName,
		})
	}
	return transactions, nil
}

func (s *transactionService) GetLastTransactions(ctx context.Context, request model.LastTransactionRequest) ([]model.TransactionResponse, error) {
	var walletIDs []int
	if request.WalletID != 0 {
		walletIDs = append(walletIDs, request.WalletID)
	} else {
		getWallets, err := s.walletRepo.GetWalletByUserID(ctx, int(request.UserID))
		if err != nil {
			return nil, err
		}
		for _, wallet := range getWallets {
			walletIDs = append(walletIDs, wallet.ID)
		}
	}

	getTransaction, err := s.transactionRepo.GetLastTransactions(ctx, walletIDs)
	if err != nil {
		return nil, err
	}

	var transactions []model.TransactionResponse
	for _, record := range getTransaction {
		var categoryID int32
		var categoryName string
		if record.CategoryID != nil {
			categoryID = int32(*record.CategoryID)
			categoryName = record.Category.Name
		} else {
			categoryID = 0
			categoryName = ""
		}
		transactions = append(transactions, model.TransactionResponse{
			ID:              int32(record.ID),
			TransactionDate: record.TransactionDate,
			Type:            record.Type,
			Nominal:         record.Nominal,
			WalletID:        int32(record.WalletID),
			WalletName:      record.Wallet.Name,
			CategoryID:      categoryID,
			CategoryName:    categoryName,
		})
	}
	return transactions, nil
}

func (s *transactionService) GetCashflowReport(ctx context.Context, request model.CashFlowReportRequest) (model.CashFlowReportResponse, error) {
	filter := model.FilterTransaction{}
	filter.StartTime = request.StartTime
	filter.EndTime = request.EndTime
	if request.WalletID != 0 {
		filter.WalletID = []int32{request.WalletID}
	} else {
		//get wallet id by userID
		getWallets, err := s.walletRepo.GetWalletByUserID(ctx, int(request.UserID))
		if err != nil {
			return model.CashFlowReportResponse{}, err
		}
		for _, wallet := range getWallets {
			filter.WalletID = append(filter.WalletID, int32(wallet.ID))
		}
	}

	result, err := s.transactionRepo.GetCashflowReport(ctx, filter)
	if err != nil {
		return model.CashFlowReportResponse{}, err
	}
	var output model.CashFlowReportResponse
	for _, record := range result {
		if record.Type == "Income" {
			output.In = record.Total
		}
		if record.Type == "Expense" {
			output.Out = record.Total
		}
	}
	return output, nil
}

func (s *transactionService) GetSummaryCategory(ctx context.Context, request model.SummaryCategoryRequest) ([]model.SummaryCategoryResponse, error) {
	filter := model.FilterTransaction{}
	filter.StartTime = request.StartTime
	filter.EndTime = request.EndTime
	if request.WalletID != 0 {
		filter.WalletID = []int32{request.WalletID}
	} else {
		//get wallet id by userID
		getWallets, err := s.walletRepo.GetWalletByUserID(ctx, int(request.UserID))
		if err != nil {
			return nil, err
		}
		for _, wallet := range getWallets {
			filter.WalletID = append(filter.WalletID, int32(wallet.ID))
		}
	}

	getSummaryCategory, err := s.transactionRepo.GetSummaryCategory(ctx, filter)
	if err != nil {
		return nil, err
	}

	var summaryCategories []model.SummaryCategoryResponse
	for _, record := range getSummaryCategory {

		summaryCategories = append(summaryCategories, model.SummaryCategoryResponse{
			CategoryID:   record.CategoryID,
			CategoryName: record.CategoryName,
			Type:         record.Type,
			Amount:       record.Amount,
		})
	}
	return summaryCategories, nil
}
