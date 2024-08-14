package model

import "time"

type SummaryCategoryRequest struct {
	UserID    int32
	WalletID  int32
	StartTime *time.Time
	EndTime   *time.Time
}
