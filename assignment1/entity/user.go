package entity

import "time"

type User struct {
	ID          uint         `gorm:"primaryKey"`
	Name        string       `gorm:"type:varchar;not null" json:"name" binding:"required"`
	Email       string       `gorm:"type:varchar;uniqueIndex;not null" json:"email" binding:"required,email"`
	CreatedAt   time.Time    `gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime"`
	Submissions []Submission `gorm:"foreignKey:UserID"`
}
