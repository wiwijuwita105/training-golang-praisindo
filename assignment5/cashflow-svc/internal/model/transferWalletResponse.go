package model

import (
	"assignment5/cashflow-svc/internal/entity"
	"time"
)

type TransferWalletResponse struct {
	SenderWalletID   int                `json:"senderWalletID"`
	ReceiverWalletID int                `json:"receiverWalletID"`
	Nominal          float64            `json:"nominal"`
	TransactionDate  time.Time          `json:"transactionDate"`
	TransactionIn    entity.Transaction `json:"transactionIn"`
	TransactionOut   entity.Transaction `json:"transactionOut"`
}
