package models

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("user not found")
	ErrUsernameConflict  = errors.New("username is already taken")
	ErrInternalServer    = errors.New("internal server error")
	ErrIncorrectPassword = errors.New("incorrect password")
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
