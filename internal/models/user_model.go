package models

import "time"

type User struct {
	Username   string
	Fullname   string
	Password   string
	Birthday   string
	Balance    float64
	Registered time.Time
	Privilege  string
}

type NewUserDTO struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
}
