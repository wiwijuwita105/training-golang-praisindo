package entity

import "time"

type TopupRequest struct {
	ID     int     `json:"ID" name:"ID"`
	Amount float64 `json:"amount" name:"amount"`
}

type TransactionResponse struct {
	Message string `json:"message"`
}

type TransactionRequest struct {
	FromID int32   `json:"fromID" name:"fromID"`
	ToID   int     `json:"toID" name:"toID"`
	Type   string  `json:"type" name:"type"`
	Amount float64 `json:"amount" name:"amount"`
}

type TransferRequest struct {
	FromID int32   `json:"fromID" name:"fromID"`
	ToID   int     `json:"toID" name:"toID"`
	Type   string  `json:"type" name:"type"`
	Amount float64 `json:"amount" name:"amount"`
}

type TransactionGetResponse struct {
	ID        int       `json:"ID" name:"ID"`
	UserID    int       `json:"userID" name:"userID"`
	Name      string    `json:"name" name:"name"`
	Amount    float64   `json:"amount" name:"amount"`
	Category  string    `json:"category" name:"category"`
	Type      string    `json:"type" name:"type"`
	CreatedAt time.Time `json:"createdAt" name:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" name:"updatedAt"`
}

type TransactionGetRequest struct {
	Type   string `json:"type" name:"type"`
	UserID int    `json:"userID" name:"userID"`
	Size   int    `json:"size" name:"size"`
	Page   int    `json:"page" name:"page"`
}

type Pagination struct {
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
	PageSize  int `json:"page_size"`
	Page      int `json:"page"`
}

type TransactionGetResponseWithPagination struct {
	Data       []TransactionGetResponse `json:"data"`
	Pagination Pagination               `json:"pagination"`
}
