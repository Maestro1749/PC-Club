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

func (s *ComputerService) DeleteComputer(id int) error {
	computer := &models.Computer{
		ID: id,
	}

	return s.repo.DeleteComputer(computer)
}

func (s *ComputerService) ChangePrice(number string, newPrice float64) error {
	if newPrice < 0 {
		return errors.New("The price should not be negative.")
	}

	return s.repo.ChangePrice(number, newPrice)
}
