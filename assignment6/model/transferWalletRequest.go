package model

type TransferWalletRequest struct {
	ToID            int32   `json:"ToID"`
	FromID          int32   `json:"FromID"`
	Nominal         float64 `json:"nominal"`
	TransactionDate string  `json:"transactionDate"`
}
