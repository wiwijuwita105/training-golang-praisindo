package entity

import "time"

type TransactionCategory struct {
	ID        int           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string        `json:"name" binding:"required"`
	Type      string        `json:"type" binding:"required"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Records   []Transaction `gorm:"foreignKey:CategoryID"`
}
