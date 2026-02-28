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

type NewComputerDTO struct {
	Num   string  `json:"number"`
	Price float64 `json:"price"`
}

type DeleteComputerDTO struct {
	ID int `json:"id"`
}

type ComputerRepository interface {
	Create(computer *Computer) error
	Delete(id int) error
}
