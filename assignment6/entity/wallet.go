package entity

import "time"

type Wallet struct {
	ID           int           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int           `json:"userID" binding:"required"`
	Name         string        `json:"name" binding:"required"`
	Balance      float64       `json:"balance"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Transactions []Transaction `gorm:"foreignKey:WalletID"`
}
