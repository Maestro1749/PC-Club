package booking_handlers

import (
	"encoding/json"
	"errors"
	"mvp/internal/models"
	"mvp/internal/service/booking_service"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type BookingHandler struct {
	service *booking_service.BookingService
	logger  *zap.Logger
}

func NewBookingHandler(service *booking_service.BookingService, logger *zap.Logger) *BookingHandler {
	return &BookingHandler{service: service, logger: logger}
}

func (h *BookingHandler) BookingComputerHandler(w http.ResponseWriter, r *http.Request) {
	var BookingDTO BookComputerDTO
	if err := json.NewDecoder(r.Body).Decode(&BookingDTO); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.BookingComputer(BookingDTO.ComputerNumber, BookingDTO.UserID, BookingDTO.StartTime, BookingDTO.DurationHours); err != nil {
		h.logger.Error("failed to book computer", zap.Error(err))

		switch {
		case errors.Is(err, models.ErrComputerNotFound):
			writeError(w, err, http.StatusNotFound)
		case errors.Is(err, models.ErrComputerBusy):
			writeError(w, err, http.StatusConflict)
		case errors.Is(err, models.ErrTimeFormat):
			writeError(w, err, http.StatusBadRequest)
		case errors.Is(err, models.ErrInvalidDuration):
			writeError(w, err, http.StatusBadRequest)
		default:
			writeError(w, err, http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

func writeError(w http.ResponseWriter, err error, status int) {
	newError := models.ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	str, _ := newError.ToString()

	http.Error(w, str, status)
}
