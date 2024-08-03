package entity

import "time"

type Transaction struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	WalletID        int       `json:"wallet_id" binding:"required"`
	CategoryID      int       `json:"category_id" binding:"required"`
	TransactionDate time.Time `json:"transaction_date" binding:"required"`
	Type            string    `json:"type" binding:"required"`
	Nominal         float64   `gorm:"type:decimal(10,4);" json:"nominal"`
	FromWalletID    int       `json:"from_wallet_id" binding:"required"`
	ToWalletID      int       `json:"to_wallet_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
