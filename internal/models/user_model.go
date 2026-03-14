package models

import (
	"time"
)

type User struct {
	ID          int
	Username    string
	Fullname    string
	Email       string
	PhoneNumber string
	Password    string
	Birthday    time.Time
	Balance     float64
	Registered  time.Time
	Privilege   string
}

type NewUserDTO struct {
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Birthday    string `json:"birthday"`
}

type LoginUserDTO struct {
	Username    string `json:"username"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}
