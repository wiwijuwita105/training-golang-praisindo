package service

import (
	"context"
	"fmt"
	"wallet_svc/entity"
)

// IWalletService mendefinisikan interface untuk layanan Wallet
type IWalletService interface {
	CreateWallet(ctx context.Context, user *entity.Wallet) (entity.Wallet, error)
	GetWalletByUserID(ctx context.Context, id int) (entity.Wallet, error)
	GetAllWallets(ctx context.Context) ([]entity.Wallet, error)
}

// IWalletRepository mendefinisikan interface untuk repository wallet
type IWalletRepository interface {
	CreateWallet(ctx context.Context, user *entity.Wallet) (entity.Wallet, error)
	GetWalletByUserID(ctx context.Context, id int) (entity.Wallet, error)
	GetAllWallets(ctx context.Context) ([]entity.Wallet, error)
}

// walletService adalah implementasi dari IWalletService yang menggunakan IWalletRepository
type walletService struct {
	walletRepo IWalletRepository
}

// NewWalletService membuat instance baru dari walletService
func NewWalletService(walletRepo IWalletRepository) IWalletService {
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
func (s *walletService) GetWalletByUserID(ctx context.Context, userid int) (entity.Wallet, error) {
	// Memanggil GeWalletByUserID dari repository untuk mendapatkan wallet berdasarkan UserID
	wallet, err := s.walletRepo.GetWalletByUserID(ctx, userid)
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
