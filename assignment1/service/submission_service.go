package service

import (
	"assignment1/entity"
	"context"
	"fmt"
)

// ISubmissionService mendefinisikan interface untuk layanan pengguna
type ISubmissionService interface {
	CreateSubmission(ctx context.Context, submission *entity.Submission) (entity.Submission, error)
	GetSubmissionByID(ctx context.Context, id int) (entity.Submission, error)
	DeleteSubmission(ctx context.Context, id int) error
	GetAllSubmissions(ctx context.Context, page int, pageSize int, userID *uint) ([]entity.Submission, int64, error)         // Menambahkan parameter page dan pageSize
	GetAllSubmissionsWithUser(ctx context.Context, page int, pageSize int, userID *uint) ([]entity.Submission, int64, error) // Menambahkan parameter page dan pageSize
}

// ISubmissionRepository mendefinisikan interface untuk repository pengguna
type ISubmissionRepository interface {
	CreateSubmission(ctx context.Context, submission *entity.Submission) (entity.Submission, error)
	GetSubmissionByID(ctx context.Context, id int) (entity.Submission, error)
	DeleteSubmission(ctx context.Context, id int) error
	GetAllSubmissions(ctx context.Context, page int, pageSize int, userID *uint) ([]entity.Submission, int64, error)         // Menambahkan parameter page dan pageSize
	GetAllSubmissionsWithUser(ctx context.Context, page int, pageSize int, userID *uint) ([]entity.Submission, int64, error) // Menambahkan parameter page dan pageSize
}

// submissionService adalah implementasi dari ISubmissionService yang menggunakan ISubmissionRepository
type submissionService struct {
	submissionRepo ISubmissionRepository
}

// NewsubmissionService membuat instance baru dari submissionService
func NewSubmissionService(submissionRepo ISubmissionRepository) ISubmissionService {
	return &submissionService{submissionRepo: submissionRepo}
}

// Create
func (s *submissionService) CreateSubmission(ctx context.Context, submission *entity.Submission) (entity.Submission, error) {
	// Memanggil CreateUser dari repository untuk membuat pengguna baru
	createdSubmission, err := s.submissionRepo.CreateSubmission(ctx, submission)
	if err != nil {
		return entity.Submission{}, fmt.Errorf("gagal membuat pengguna: %v", err)
	}
	return createdSubmission, nil
}

// GetUserByID mendapatkan pengguna berdasarkan ID
func (s *submissionService) GetSubmissionByID(ctx context.Context, id int) (entity.Submission, error) {
	// Memanggil GetUserByID dari repository untuk mendapatkan pengguna berdasarkan ID
	submission, err := s.submissionRepo.GetSubmissionByID(ctx, id)
	if err != nil {
		return entity.Submission{}, fmt.Errorf("gagal mendapatkan submission berdasarkan ID: %v", err)
	}
	return submission, nil
}

// Delete submission
func (s *submissionService) DeleteSubmission(ctx context.Context, id int) error {
	// Memanggil DeleteUser dari repository untuk menghapus pengguna berdasarkan ID
	err := s.submissionRepo.DeleteSubmission(ctx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus submission: %v", err)
	}
	return nil
}

// GetAllSubmission
func (s *submissionService) GetAllSubmissions(ctx context.Context, page int, pageSize int, userID *uint) ([]entity.Submission, int64, error) {
	// Memanggil GetAllUsers dari repository untuk mendapatkan semua pengguna
	Submissions, total, err := s.submissionRepo.GetAllSubmissions(ctx, page, pageSize, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mendapatkan semua submissions: %v", err)
	}
	return Submissions, total, nil
}

func (s *submissionService) GetAllSubmissionsWithUser(ctx context.Context, page int, pageSize int, userID *uint) ([]entity.Submission, int64, error) {
	// Memanggil GetAllUsers dari repository untuk mendapatkan semua pengguna
	Submissions, total, err := s.submissionRepo.GetAllSubmissionsWithUser(ctx, page, pageSize, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mendapatkan semua submissions: %v", err)
	}
	return Submissions, total, nil
}
