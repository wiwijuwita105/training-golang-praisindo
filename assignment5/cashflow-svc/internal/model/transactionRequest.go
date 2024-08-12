package model

type TransactionRequest struct {
	WalletID   int32
	CategoryID int32
	Nominal    float64
}
