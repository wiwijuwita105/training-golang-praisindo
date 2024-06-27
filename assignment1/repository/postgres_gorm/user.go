package postgres_gorm

import (
	"assignment1/entity"
	"assignment1/service"
	"context"
	"errors"
	"log"

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
	Count(count *int64) *gorm.DB
	Table(string, ...interface{}) *gorm.DB
}

type userRepository struct {
	db GormDBIface
}

// NewUserRepository membuat instance baru dari userRepository
func NewUserRepository(db GormDBIface) service.IUserRepository {
	return &userRepository{db: db}
}

// CreateUser membuat pengguna baru dalam basis data
func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("Error creating user: %v\n", err)
		return entity.User{}, err
	}
	return *user, nil
}

// GetUserByID mengambil pengguna berdasarkan ID
func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Preload("Submissions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Limit(1)
	}).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		log.Printf("Error getting user by ID: %v\n", err)
		return entity.User{}, err
	}
	return user, nil
}

// UpdateUser memperbarui informasi pengguna dalam basis data
func (r *userRepository) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	// Menemukan pengguna yang akan diperbarui
	var existingUser entity.User
	if err := r.db.WithContext(ctx).Select("id", "name", "email", "created_at", "updated_at").First(&existingUser, id).Error; err != nil {
		log.Printf("Error finding user to update: %v\n", err)
		return entity.User{}, err
	}

	// Memperbarui informasi pengguna
	existingUser.Name = user.Name
	existingUser.Email = user.Email
	if err := r.db.WithContext(ctx).Save(&existingUser).Error; err != nil {
		log.Printf("Error updating user: %v\n", err)
		return entity.User{}, err
	}
	return existingUser, nil
}

// DeleteUser menghapus pengguna berdasarkan ID
func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error; err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil
}

// GetAllUsers mengambil semua pengguna dari basis data
func (r *userRepository) GetAllUsers(ctx context.Context, page int, pageSize int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	offset := (page - 1) * pageSize

	//count data
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&total).Error; err != nil {
		log.Printf("Error counting users: %v\n", err)
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Select("id", "name", "email", "created_at", "updated_at").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, total, nil
		}
		log.Printf("Error getting all users: %v\n", err)
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) GetAllUsersWithRisk(ctx context.Context, page int, pageSize int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	offset := (page - 1) * pageSize

	//count data
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&total).Error; err != nil {
		log.Printf("Error counting users: %v\n", err)
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Submissions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Limit(1)
	}).Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, total, nil
		}
		log.Printf("Error getting all users: %v\n", err)
		return nil, 0, err
	}
	return users, total, nil
}
