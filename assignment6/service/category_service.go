package service

import (
	"assignment6/entity"
	"assignment6/repository/postgres_gorm"
	"context"
	"fmt"
)

// ICategoryService mendefinisikan interface untuk layanan Wallet
type ICategoryService interface {
	CreateCategory(ctx context.Context, category *entity.TransactionCategory) (entity.TransactionCategory, error)
	GetCategoryByID(ctx context.Context, id int) (entity.TransactionCategory, error)
	GetAllCategories(ctx context.Context) ([]entity.TransactionCategory, error)
	UpdateCategory(ctx context.Context, id int, category entity.TransactionCategory) (entity.TransactionCategory, error)
	DeleteCategory(ctx context.Context, id int) error
}

// categoryService adalah implementasi dari IcategoryService yang menggunakan IcategoryRepository
type categoryService struct {
	categoryRepo postgres_gorm.ICategoryRepository
}

// NewcategoryService membuat instance baru dari categoryService
func NewCategoryService(categoryRepo postgres_gorm.ICategoryRepository) ICategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

// Create Wallet
func (s *categoryService) CreateCategory(ctx context.Context, category *entity.TransactionCategory) (entity.TransactionCategory, error) {
	// Memanggil CreateCategory dari repository untuk membuat wallet baru
	createdCategory, err := s.categoryRepo.CreateCategory(ctx, category)
	if err != nil {
		return entity.TransactionCategory{}, fmt.Errorf("gagal membuat category: %v", err)
	}
	return createdCategory, nil
}

// GetWalletByUserID mendapatkan wallet berdasarkan User ID
func (s *categoryService) GetCategoryByID(ctx context.Context, userid int) (entity.TransactionCategory, error) {
	// Memanggil GeWalletByUserID dari repository untuk mendapatkan wallet berdasarkan UserID
	category, err := s.categoryRepo.GetCategoryByID(ctx, userid)
	if err != nil {
		return entity.TransactionCategory{}, fmt.Errorf("gagal mendapatkan wallet berdasarkan User ID: %v", err)
	}
	return category, nil
}

// GetAllWallets mendapatkan semua wallets
func (s *categoryService) GetAllCategories(ctx context.Context) ([]entity.TransactionCategory, error) {
	// Memanggil GetAllWallets dari repository untuk mendapatkan semua wallet
	categories, err := s.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan semua kategori: %v", err)
	}
	return categories, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id int, category entity.TransactionCategory) (entity.TransactionCategory, error) {
	updatedCategory, err := s.categoryRepo.UpdateCategory(ctx, id, category)
	if err != nil {
		return entity.TransactionCategory{}, fmt.Errorf("gagal memperbarui category: %v", err)
	}
	return updatedCategory, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id int) error {
	// Memanggil DeleteUser dari repository untuk menghapus pengguna berdasarkan ID
	err := s.categoryRepo.DeleteCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus category: %v", err)
	}
	return nil
}
