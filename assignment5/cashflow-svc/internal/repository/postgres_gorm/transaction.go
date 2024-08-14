package postgres_gorm

import (
	"assignment5/cashflow-svc/internal/entity"
	"assignment5/cashflow-svc/internal/model"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db GormDBIface
}

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error)
	GetTransactionByID(ctx context.Context, id int) (entity.Transaction, error)
	GetAllTransactions(ctx context.Context, filter model.FilterTransaction) ([]entity.Transaction, error)
	DeleteTransaction(ctx context.Context, id int) error
	GetLastTransactions(ctx context.Context, walletID []int) ([]entity.Transaction, error)
	GetCashflowReport(ctx context.Context, filter model.FilterTransaction) ([]model.CashFlowResult, error)
	GetSummaryCategory(ctx context.Context, filte model.FilterTransaction) ([]model.SummaryCategoryResponse, error)
}

func NewTransactionRepository(db GormDBIface) ITransactionRepository {

	return &transactionRepository{db: db}
}

// CreateWallet membuat pengguna baru dalam basis data
func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error) {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		log.Printf("Error creating user: %v\n", err)
		return entity.Transaction{}, err
	}
	return *transaction, nil
}

// GetWalletByUserID mengambil wallet berdasarkan User ID
func (r *transactionRepository) GetTransactionByID(ctx context.Context, id int) (entity.Transaction, error) {
	var transaction entity.Transaction
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&transaction).Error; err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Transaction{}, nil
		}
		log.Printf("Error getting wallet by user ID: %v\n", err)
		return entity.Transaction{}, err
	}
	log.Println(transaction)
	return transaction, nil
}

// GetAllWalletss mengambil semua wallet dari basis data
func (r *transactionRepository) GetAllTransactions(ctx context.Context, filter model.FilterTransaction) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.WithContext(ctx).
		Where("wallet_id IN ?", filter.WalletID).
		Where("transaction_date BETWEEN ? AND ?", filter.StartTime, filter.EndTime).
		Preload("Wallet").
		Preload("Category").
		Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactions, nil
		}
		log.Printf("Error getting all wallets: %v\n", err)
		return nil, err
	}
	//log.Fatalln(transactions)
	return transactions, nil
}

func (r *transactionRepository) GetLastTransactions(ctx context.Context, walletIDs []int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.WithContext(ctx).
		Where("wallet_id IN ?", walletIDs).
		Preload("Wallet").
		Preload("Category").
		Order("transaction_date desc").
		Limit(10).
		Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactions, nil
		}
		log.Printf("Error getting all wallets: %v\n", err)
		return nil, err
	}
	//log.Fatalln(transactions)
	return transactions, nil
}

func (r *transactionRepository) DeleteTransaction(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Transaction{}, id).Error; err != nil {
		log.Printf("Error deleting Transaction: %v\n", err)
		return err
	}
	return nil
}

func (r *transactionRepository) GetCashflowReport(ctx context.Context, filter model.FilterTransaction) ([]model.CashFlowResult, error) {
	var results []model.CashFlowResult

	// Perform a grouped query to sum amounts by type
	err := r.db.WithContext(ctx).Model(&entity.Transaction{}).
		Select("type, sum(nominal) as total").
		Where("transaction_date BETWEEN ? AND ?", filter.StartTime, filter.EndTime).
		Where("wallet_id IN ?", filter.WalletID).
		Group("type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *transactionRepository) GetSummaryCategory(ctx context.Context, filter model.FilterTransaction) ([]model.SummaryCategoryResponse, error) {
	var results []model.SummaryCategoryResponse

	err := r.db.WithContext(ctx).Model(&entity.Transaction{}).
		Table("transactions as t"). // Ensure the table alias is used
		Select("t.category_id, c.name as category_name, t.type, SUM(t.nominal) as amount").
		Joins("JOIN transaction_categories as c on c.id = t.category_id"). // Correct join, assuming table name and foreign key
		Where("t.transaction_date BETWEEN ? AND ?", filter.StartTime, filter.EndTime).
		Where("t.wallet_id IN ?", filter.WalletID).
		Where("t.category_id IS NOT NULL").
		Group("t.category_id, c.name, t.type"). // Include all non-aggregate columns here
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
