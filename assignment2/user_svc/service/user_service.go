package service

import (
	"context"
	"fmt"
	"user_svc/entity"
)

// IUserService mendefinisikan interface untuk layanan pengguna
type IUserService interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

// IUserRepository mendefinisikan interface untuk repository pengguna
type IUserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
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

// GetAllUsers mendapatkan semua pengguna
func (s *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	// Memanggil GetAllUsers dari repository untuk mendapatkan semua pengguna
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan semua pengguna: %v", err)
	}
	return users, nil
}
