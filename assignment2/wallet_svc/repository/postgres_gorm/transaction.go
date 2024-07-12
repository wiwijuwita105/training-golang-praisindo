package postgres_gorm

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"wallet_svc/entity"
	"wallet_svc/service"
)

type transactionRepository struct {
	db GormDBIface
}

func NewTransactionRepository(db GormDBIface) service.ITransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error) {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		log.Printf("Error creating transaction: %v\n", err)
		return entity.Transaction{}, err
	}
	return *transaction, nil
}

func (r *transactionRepository) GetAllTransactions(ctx context.Context, param entity.TransactionGetRequest) ([]entity.Transaction, error) {
	log.Println(param.Type)
	var transactions []entity.Transaction
	if err := r.db.WithContext(ctx).Where("type = ?", param.Type).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactions, nil
		}
		log.Printf("Error getting all wallets: %v\n", err)
		return nil, err
	}
	return transactions, nil
}
