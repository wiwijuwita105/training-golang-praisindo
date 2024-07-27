package entity

import "time"

type Shortlink struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Shortlink string    `gorm:"type:varchar;uniqueIndex;not null" json:"shortlink" binding:"required"`
	Url       string    `gorm:"type:varchar;not null" json:"url" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
