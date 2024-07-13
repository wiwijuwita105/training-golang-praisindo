package entity

import "time"

type Transaction struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    int       `json:"userID" binding:"required"`
	Amount    float64   `json:"amount" binding:"required"`
	Category  string    `gorm:"not null" validate:"required,oneof=in out"`
	Type      string    `gorm:"not null" binding:"required,oneof=topup transfer"`
	Timestamp time.Time `gorm:"autoCreateTime" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransactionRequest struct {
	FromID int     `json:"fromID" name:"fromID"`
	ToID   int     `json:"toID" name:"toID"`
	Type   string  `json:"type" name:"type"`
	Amount float64 `json:"amount" name:"amount"`
}

type TransactionGetRequest struct {
	Type   string `json:"type" name:"type"`
	UserID int    `json:"userID" name:"userID"`
	Size   int    `json:"size" name:"size"`
	Page   int    `json:"page" name:"page"`
}

type TransactionResponse struct {
	Transaction []Transaction `json:"transaction"`
	CountData   int64         `json:"countData"`
}

type Pagination struct {
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
	PageSize  int `json:"page_size"`
	Page      int `json:"page"`
}

type TransactionGetResponseWithPagination struct {
	Data       []Transaction `json:"data"`
	Pagination Pagination    `json:"pagination"`
}
