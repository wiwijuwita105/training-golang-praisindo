package entity

import "time"

type Wallet struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    int       `json:"userID" binding:"required"`
	Balance   float64   `json:"balance" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
