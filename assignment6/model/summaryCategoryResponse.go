package model

type SummaryCategoryResponse struct {
	CategoryID   int32   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
}
