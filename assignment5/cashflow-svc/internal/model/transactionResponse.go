package model

import "time"

type TransactionResponse struct {
	ID              int32     `json:"id"`
	TransactionDate time.Time `json:"transaction_date"`
	Nominal         float64   `json:"nominal"`
	Type            string    `json:"type"`
	WalletID        int32     `json:"walletId"`
	WalletName      string    `json:"walletName"`
	CategoryID      int32     `json:"categoryId"`
	CategoryName    string    `json:"categoryName"`
}
