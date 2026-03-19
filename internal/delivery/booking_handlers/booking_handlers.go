package booking_handlers

import (
	"encoding/json"
	"errors"
	"mvp/internal/models"
	"mvp/internal/service/booking_service"
	"net/http"

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

		if errors.Is(err, models.ErrComputerNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if errors.Is(err, models.ErrComputerBusy) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		if errors.Is(err, models.ErrTimeFormat) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, models.ErrInvalidDuration) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
