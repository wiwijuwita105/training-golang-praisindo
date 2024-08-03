package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Fullname  string    `gorm:"type:varchar;not null" json:"name" binding:"required"`
	Email     string    `gorm:"type:varchar;uniqueIndex;not null" json:"email" binding:"required,email"`
	Password  string    `gorm:"type:varchar;not null" json:"password" binding:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Wallets   []Wallet  `gorm:"foreignKey:UserID"`
}
