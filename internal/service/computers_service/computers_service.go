package computers_service

import (
	"mvp/internal/models"
	"mvp/internal/repository/computers_repository"

	"go.uber.org/zap"
)

type ComputerService struct {
	repo   computers_repository.ComputerRepository
	logger *zap.Logger
}

func NewComputerService(repo computers_repository.ComputerRepository, logger *zap.Logger) *ComputerService {
	return &ComputerService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ComputerService) AddComputer(number string, price float64) error {
	if price < 0 {
		s.logger.Error("Invalid price for computer", zap.Float64("price", price))
		return models.ErrPriceConflict
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
		s.logger.Error("Invalid new price for computer", zap.Float64("new price", newPrice))
		return models.ErrPriceConflict
	}

	return s.repo.ChangePrice(number, newPrice)
}
