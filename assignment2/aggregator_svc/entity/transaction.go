package entity

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
