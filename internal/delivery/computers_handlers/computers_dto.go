package computers_handlers

type NewComputerDTO struct {
	Num   string  `json:"number"`
	Price float64 `json:"price"`
}

type DeleteComputerDTO struct {
	ID int `json:"id"`
}

type ChangePriceDTO struct {
	Number   string  `json:"number"`
	NewPrice float64 `json:"newPrice"`
}
