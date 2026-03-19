package booking_service

import (
	"errors"
	"mvp/internal/models"
	"mvp/internal/repository/booking_repository"
	"time"

	"go.uber.org/zap"
)

type BookingService struct {
	repo   booking_repository.BookingRepository
	logger *zap.Logger
}

func NewBookingService(repo booking_repository.BookingRepository, logger *zap.Logger) *BookingService {
	return &BookingService{
		repo:   repo,
		logger: logger,
	}
}

func (s *BookingService) BookingComputer(computerNum string, userID int, timeStart string, durationsHours float64) error {
	if durationsHours < 0.5 || durationsHours > 24 {
		return models.ErrInvalidDuration
	}

	exist, err := s.repo.Exists(computerNum)
	if err != nil {
		return err
	}
	if !exist {
		return models.ErrComputerNotFound
	}

	id, err := s.repo.FindComputerIDByNumber(computerNum)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		s.logger.Error("error to load location")
		return models.ErrInternalServer
	}

	start, err := time.ParseInLocation("02.01.2006 15:04", timeStart, loc)
	if err != nil {
		return models.ErrTimeFormat
	}

	if time := time.Until(start); time < 0 {
		return models.ErrTimeFormat
	}

	timeEnd := start.Add(time.Duration(durationsHours * float64(time.Hour)))

	isBusy, err := s.repo.BusyCheck(start, timeEnd, id)
	if err != nil {
		return err
	}
	if isBusy {
		return errors.New("computer is not available")
	}

	price, err := s.repo.TakePrice(computerNum)
	if err != nil {
		return err
	}

	price *= durationsHours

	if err := s.repo.Booking(userID, price, id, start, timeEnd); err != nil {
		return err
	}

	return nil
}
