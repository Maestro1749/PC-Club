package models

import "time"

type Computer struct {
	Number    string
	Price     float64
	IsBusy    bool
	BusySince *time.Time
	BusyUntil *time.Time
}

type NewComputerDTO struct {
	Number string
	Price  float64
}

type DeleteComputerDTO struct {
	Number string
}
