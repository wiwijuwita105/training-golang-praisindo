package postgres_gorm

import (
	"assignment5/cashflow-svc/internal/entity"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type walletRepository struct {
	db GormDBIface
}

// IWalletRepository mendefinisikan interface untuk repository wallet
type IWalletRepository interface {
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error)
	GetWalletByID(ctx context.Context, id int) (entity.Wallet, error)
	GetAllWallets(ctx context.Context) ([]entity.Wallet, error)
	UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error)
	DeleteWallet(ctx context.Context, id int) error
	GetWalletByUserID(ctx context.Context, userID int) ([]entity.Wallet, error)
}

// NewWalletRepository membuat instance baru dari walletRepository
func NewWalletRepository(db GormDBIface) IWalletRepository {
	return &walletRepository{db: db}
}

// CreateWallet membuat pengguna baru dalam basis data
func (r *walletRepository) CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error) {
	if err := r.db.WithContext(ctx).Create(wallet).Error; err != nil {
		log.Printf("Error creating wallet: %v\n", err)
		return entity.Wallet{}, err
	}
	return *wallet, nil
}

// GetWalletByUserID mengambil wallet berdasarkan User ID
func (r *walletRepository) GetWalletByID(ctx context.Context, id int) (entity.Wallet, error) {
	var wallet entity.Wallet
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&wallet).Error; err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Wallet{}, nil
		}
		log.Printf("Error getting wallet by user ID: %v\n", err)
		return entity.Wallet{}, err
	}
	log.Println(wallet)
	return wallet, nil
}

// GetAllWalletss mengambil semua wallet dari basis data
func (r *walletRepository) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
	var wallets []entity.Wallet
	if err := r.db.WithContext(ctx).Find(&wallets).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return wallets, nil
		}
		log.Printf("Error getting all wallets: %v\n", err)
		return nil, err
	}
	return wallets, nil
}

func (r *walletRepository) UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error) {
	// Menemukan pengguna yang akan diperbarui
	var existingWallet entity.Wallet
	if err := r.db.WithContext(ctx).First(&existingWallet, id).Error; err != nil {
		log.Printf("Error finding wallet to update: %v\n", err)
		return entity.Wallet{}, err
	}

	existingWallet.Balance = wallet.Balance
	if err := r.db.WithContext(ctx).Save(&existingWallet).Error; err != nil {
		log.Printf("Error updating wallet: %v\n", err)
		return entity.Wallet{}, err
	}
	return existingWallet, nil
}

func (r *walletRepository) DeleteWallet(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Wallet{}, id).Error; err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil
}

func (r *walletRepository) GetWalletByUserID(ctx context.Context, userID int) ([]entity.Wallet, error) {
	var wallets []entity.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&wallets).Error; err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Wallet{}, nil
		}
		log.Printf("Error getting wallet by user ID: %v\n", err)
		return []entity.Wallet{}, err
	}

	return wallets, nil
}
