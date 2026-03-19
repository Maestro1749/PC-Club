package models

import "time"

type Booking struct {
	ID         int
	UserID     int
	ComputerID int
	StartTime  time.Time
	EndTime    time.Time
	TotalCost  float64
	CreatedAt  time.Time
}
