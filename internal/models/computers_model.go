package models

import "time"

type Computer struct {
	ID        int
	Num       string
	Price     float64
	IsBusy    bool
	BusySince *time.Time
	BusyUntil *time.Time
}
