package models

import (
	"encoding/json"
	"time"
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
