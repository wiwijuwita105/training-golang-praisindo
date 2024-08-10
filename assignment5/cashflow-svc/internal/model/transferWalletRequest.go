package model

import "time"

type TransferWalletRequest struct {
	ToID            int32     `json:"ToID"`
	FromID          int32     `json:"FromID"`
	Nominal         float64   `json:"nominal"`
	TransactionDate time.Time `json:"transactionDate"`
}
