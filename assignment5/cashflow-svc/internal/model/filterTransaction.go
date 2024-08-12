package model

import "time"

type FilterTransaction struct {
	StartTime *time.Time
	EndTime   *time.Time
	WalletID  []int32
	UserID    int32
}
