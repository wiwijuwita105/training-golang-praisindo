package postgres_gorm

import (
	"assignment1/entity"
	"assignment1/service"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type submissionRepository struct {
	db GormDBIface
}

// NewUserRepository membuat instance baru dari userRepository
func NewSubmissionRepository(db GormDBIface) service.ISubmissionRepository {
	return &submissionRepository{db: db}
}

// Create submission
func (r *submissionRepository) CreateSubmission(ctx context.Context, submission *entity.Submission) (entity.Submission, error) {
	if err := r.db.WithContext(ctx).Create(submission).Error; err != nil {
		log.Printf("Error creating submission: %v\n", err)
		return entity.Submission{}, err
	}
	return *submission, nil
}

// GetUserByID mengambil pengguna berdasarkan ID
func (r *submissionRepository) GetSubmissionByID(ctx context.Context, id int) (entity.Submission, error) {
	var submission entity.Submission
	if err := r.db.WithContext(ctx).First(&submission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Submission{}, nil
		}
		log.Printf("Error getting submission by ID: %v\n", err)
		return entity.Submission{}, err
	}
	return submission, nil
}

// Deletesubmission menghapus pengguna berdasarkan ID
func (r *submissionRepository) DeleteSubmission(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Submission{}, id).Error; err != nil {
		log.Printf("Error deleting submission: %v\n", err)
		return err
	}
	return nil
}

// GetAll Submission
func (r *submissionRepository) GetAllSubmissions(ctx context.Context, page int, pageSize int, userID *uint) ([]entity.Submission, int64, error) {
	var submissions []entity.Submission
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&entity.Submission{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	//query count data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(pageSize).Offset(offset).Find(&submissions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return submissions, total, nil
		}
		log.Printf("Error getting all submissions: %v\n", err)
		return nil, 0, err
	}

	return submissions, total, nil
}

func (r *submissionRepository) GetAllSubmissionsWithUser(ctx context.Context, page int, pageSize int, userID *uint) ([]entity.Submission, int64, error) {
	var submissions []entity.Submission
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&entity.Submission{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	//query count data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Preload("User").Limit(pageSize).Order("created_at DESC").Offset(offset).Find(&submissions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return submissions, total, nil
		}
		log.Printf("Error getting all submissions: %v\n", err)
		return nil, 0, err
	}

	return submissions, total, nil
}
