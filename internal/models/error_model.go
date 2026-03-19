package models

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	// user errors
	ErrNotFound            = errors.New("user not found")
	ErrUsernameConflict    = errors.New("username is already taken")
	ErrEmailConflict       = errors.New("email is already taken")
	ErrPhoneNumberConflict = errors.New("phone number is already taken")
	ErrWrongPassword       = errors.New("wrong password")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidUsername     = errors.New("invalid username")
	ErrInvalidBirtday      = errors.New("invalid birthday")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrInvalidPhoneNumber  = errors.New("invalid phone number")

	// computer errors
	ErrComputerNumConflict = errors.New("num is already taken")
	ErrComputerNotFound    = errors.New("computer not found")
	ErrPriceConflict       = errors.New("the price should not be negative")

	// booking errors
	ErrTimeFormat      = errors.New("invalid time format")
	ErrComputerBusy    = errors.New("computer is not available")
	ErrInvalidDuration = errors.New("invalid duration")

	// internal errors
	ErrInternalServer = errors.New("internal server error")
)

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e ErrorDTO) ToString() (string, error) {
	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
