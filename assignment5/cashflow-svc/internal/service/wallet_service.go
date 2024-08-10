package service

import (
	"assignment5/cashflow-svc/internal/entity"
	"assignment5/cashflow-svc/internal/repository/postgres_gorm"
	"context"
	"fmt"
)

// IWalletService mendefinisikan interface untuk layanan Wallet
type IWalletService interface {
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error)
	GetWalletByID(ctx context.Context, id int) (entity.Wallet, error)
	GetAllWallets(ctx context.Context) ([]entity.Wallet, error)
	UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error)
	DeleteWallet(ctx context.Context, id int) error
}

// walletService adalah implementasi dari IWalletService yang menggunakan IWalletRepository
type walletService struct {
	walletRepo postgres_gorm.IWalletRepository
}

// NewWalletService membuat instance baru dari walletService
func NewWalletService(walletRepo postgres_gorm.IWalletRepository) IWalletService {
	return &walletService{walletRepo: walletRepo}
}

// Create Wallet
func (s *walletService) CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error) {
	// Memanggil CreateWallet dari repository untuk membuat wallet baru
	createdWallet, err := s.walletRepo.CreateWallet(ctx, wallet)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("gagal membuat wallet: %v", err)
	}
	return createdWallet, nil
}

// GetWalletByUserID mendapatkan wallet berdasarkan User ID
func (s *walletService) GetWalletByID(ctx context.Context, userid int) (entity.Wallet, error) {
	// Memanggil GeWalletByUserID dari repository untuk mendapatkan wallet berdasarkan UserID
	wallet, err := s.walletRepo.GetWalletByID(ctx, userid)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("gagal mendapatkan wallet berdasarkan User ID: %v", err)
	}
	return wallet, nil
}

// GetAllWallets mendapatkan semua wallets
func (s *walletService) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
	// Memanggil GetAllWallets dari repository untuk mendapatkan semua wallet
	wallets, err := s.walletRepo.GetAllWallets(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan semua wallet: %v", err)
	}
	return wallets, nil
}

func (s *walletService) UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error) {
	updatedWallet, err := s.walletRepo.UpdateWallet(ctx, id, wallet)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("gagal memperbarui wallet: %v", err)
	}
	return updatedWallet, nil
}

func (s *walletService) DeleteWallet(ctx context.Context, id int) error {
	// Memanggil DeleteUser dari repository untuk menghapus pengguna berdasarkan ID
	err := s.walletRepo.DeleteWallet(ctx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus wallet: %v", err)
	}
	return nil
}
