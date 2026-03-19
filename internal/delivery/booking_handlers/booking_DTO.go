package booking_handlers

type BookComputerDTO struct {
	ComputerNumber string  `json:"computer_number"`
	UserID         int     `json:"user_id"`
	StartTime      string  `json:"start_time"`
	DurationHours  float64 `json:"duration_hours"`
}
