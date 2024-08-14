package model

type LastTransactionRequest struct {
	WalletID int `json:"walletID"`
	UserID   int `json:"userID"`
}
