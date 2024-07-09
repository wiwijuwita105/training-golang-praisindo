package postgres_gorm

import (
	"context"
	"errors"
	"log"
	"wallet_svc/entity"
	"wallet_svc/service"

	"gorm.io/gorm"
)

// GormDBIface defines an interface for GORM DB methods used in the repository
type GormDBIface interface {
	WithContext(ctx context.Context) *gorm.DB
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
}

type walletRepository struct {
	db GormDBIface
}

// NewWalletRepository membuat instance baru dari walletRepository
func NewWalletRepository(db GormDBIface) service.IWalletRepository {
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
func (r *walletRepository) GetWalletByUserID(ctx context.Context, userid int) (entity.Wallet, error) {
	var wallet entity.Wallet
	log.Println(userid)
	if err := r.db.WithContext(ctx).Where("user_id = ?", userid).Error; err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Wallet{}, nil
		}
		log.Printf("Error getting wallet by user ID: %v\n", err)
		return entity.Wallet{}, err
	}
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
