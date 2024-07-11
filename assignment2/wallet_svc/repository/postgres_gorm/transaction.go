package postgres_gorm

import (
	"context"
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
