package computers_service

import (
	"errors"
	"mvp/internal/models"
	"mvp/internal/repository/computers_repository"
)

type ComputerService struct {
	repo computers_repository.ComputerRepository
}

func NewComputerService(repo computers_repository.ComputerRepository) *ComputerService {
	return &ComputerService{repo: repo}
}

func (s *ComputerService) AddComputer(number string, price float64) error {
	if price < 0 {
		return errors.New("The price should not be negative.")
	}

	computer := &models.Computer{
		Num:   number,
		Price: price,
	}

	return s.repo.CreateComputer(computer)
}
