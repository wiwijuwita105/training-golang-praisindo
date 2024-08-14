package model

import "time"

type CashFlowReportRequest struct {
	UserID    int32
	WalletID  int32
	StartTime *time.Time
	EndTime   *time.Time
}
