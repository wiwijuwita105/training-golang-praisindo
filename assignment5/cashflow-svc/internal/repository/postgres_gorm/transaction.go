package postgres_gorm

import (
	"assignment5/cashflow-svc/internal/entity"
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
	GetAllTransactions(ctx context.Context) ([]entity.Transaction, error)
	DeleteTransaction(ctx context.Context, id int) error
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
func (r *transactionRepository) GetAllTransactions(ctx context.Context) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.WithContext(ctx).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactions, nil
		}
		log.Printf("Error getting all wallets: %v\n", err)
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) DeleteTransaction(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Transaction{}, id).Error; err != nil {
		log.Printf("Error deleting Transaction: %v\n", err)
		return err
	}
	return nil
}
