package postgres_gorm

import (
	"assignment5/cashflow-svc/internal/entity"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db GormDBIface
}

// IcategoryRepository mendefinisikan interface untuk repository wallet
type ICategoryRepository interface {
	CreateCategory(ctx context.Context, category *entity.TransactionCategory) (entity.TransactionCategory, error)
	GetCategoryByID(ctx context.Context, id int) (entity.TransactionCategory, error)
	GetAllCategories(ctx context.Context) ([]entity.TransactionCategory, error)
	UpdateCategory(ctx context.Context, id int, category entity.TransactionCategory) (entity.TransactionCategory, error)
	DeleteCategory(ctx context.Context, id int) error
}

// NewcategoryRepository membuat instance baru dari categoryRepository
func NewCategoryRepository(db GormDBIface) ICategoryRepository {
	return &categoryRepository{db: db}
}

// CreateWallet membuat pengguna baru dalam basis data
func (r *categoryRepository) CreateCategory(ctx context.Context, category *entity.TransactionCategory) (entity.TransactionCategory, error) {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		log.Printf("Error creating wallet: %v\n", err)
		return entity.TransactionCategory{}, err
	}
	return *category, nil
}

// GetWalletByUserID mengambil wallet berdasarkan User ID
func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.TransactionCategory, error) {
	var category entity.TransactionCategory
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.TransactionCategory{}, nil
		}
		log.Printf("Error getting wallet by user ID: %v\n", err)
		return entity.TransactionCategory{}, err
	}
	log.Println(category)
	return category, nil
}

// GetAllCategoriess mengambil semua wallet dari basis data
func (r *categoryRepository) GetAllCategories(ctx context.Context) ([]entity.TransactionCategory, error) {
	var categories []entity.TransactionCategory
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return categories, nil
		}
		log.Printf("Error getting all categories: %v\n", err)
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, id int, category entity.TransactionCategory) (entity.TransactionCategory, error) {
	// Menemukan pengguna yang akan diperbarui
	var existingCategory entity.TransactionCategory
	if err := r.db.WithContext(ctx).First(&existingCategory, id).Error; err != nil {
		log.Printf("Error finding wallet to update: %v\n", err)
		return entity.TransactionCategory{}, err
	}
	//tambahin value baru

	if err := r.db.WithContext(ctx).Save(&existingCategory).Error; err != nil {
		log.Printf("Error updating wallet: %v\n", err)
		return entity.TransactionCategory{}, err
	}
	return existingCategory, nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&entity.TransactionCategory{}, id).Error; err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil
}
