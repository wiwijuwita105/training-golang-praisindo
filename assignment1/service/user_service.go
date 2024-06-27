package service

import (
	"assignment1/entity"
	"context"
	"fmt"
)

// IUserService mendefinisikan interface untuk layanan pengguna
type IUserService interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context, page int, pageSize int) ([]entity.User, int64, error)         // Menambahkan parameter page dan pageSize
	GetAllUsersWithRisk(ctx context.Context, page int, pageSize int) ([]entity.User, int64, error) // Menambahkan parameter page dan pageSize
}

// IUserRepository mendefinisikan interface untuk repository pengguna
type IUserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context, page int, pageSize int) ([]entity.User, int64, error)         // Menambahkan parameter page dan pageSize
	GetAllUsersWithRisk(ctx context.Context, page int, pageSize int) ([]entity.User, int64, error) // Menambahkan parameter page dan pageSize
}

// userService adalah implementasi dari IUserService yang menggunakan IUserRepository
type userService struct {
	userRepo IUserRepository
}

// NewUserService membuat instance baru dari userService
func NewUserService(userRepo IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

// CreateUser membuat pengguna baru
func (s *userService) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	// Memanggil CreateUser dari repository untuk membuat pengguna baru
	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("gagal membuat pengguna: %v", err)
	}
	return createdUser, nil
}

// GetUserByID mendapatkan pengguna berdasarkan ID
func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	// Memanggil GetUserByID dari repository untuk mendapatkan pengguna berdasarkan ID
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("gagal mendapatkan pengguna berdasarkan ID: %v", err)
	}
	return user, nil
}

// UpdateUser memperbarui data pengguna
func (s *userService) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	// Memanggil UpdateUser dari repository untuk memperbarui data pengguna
	updatedUser, err := s.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("gagal memperbarui pengguna: %v", err)
	}
	return updatedUser, nil
}

// DeleteUser menghapus pengguna berdasarkan ID
func (s *userService) DeleteUser(ctx context.Context, id int) error {
	// Memanggil DeleteUser dari repository untuk menghapus pengguna berdasarkan ID
	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus pengguna: %v", err)
	}
	return nil
}

// GetAllUsers mendapatkan semua pengguna
func (s *userService) GetAllUsers(ctx context.Context, page int, pageSize int) ([]entity.User, int64, error) {
	// Memanggil GetAllUsers dari repository untuk mendapatkan semua pengguna
	users, total, err := s.userRepo.GetAllUsers(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mendapatkan semua pengguna: %v", err)
	}
	return users, total, nil
}

func (s *userService) GetAllUsersWithRisk(ctx context.Context, page int, pageSize int) ([]entity.User, int64, error) {
	users, total, err := s.userRepo.GetAllUsersWithRisk(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mendapatkan semua pengguna: %v", err)
	}
	return users, total, nil
}
