package computers_handlers

import (
	"encoding/json"
	"errors"
	"mvp/internal/models"
	"mvp/internal/service/computers_service"
	"net/http"
	"time"
)

type Handler struct {
	ComputerService *computers_service.ComputerService
}

func NewHandler(computerService *computers_service.ComputerService) *Handler {
	return &Handler{ComputerService: computerService}
}

func (h *Handler) AddComputerHandler(w http.ResponseWriter, r *http.Request) {
	var newComputerDTO NewComputerDTO
	if err := json.NewDecoder(r.Body).Decode(&newComputerDTO); err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, "Error: Incorrect data struct.", http.StatusBadRequest)
			return
		}

		http.Error(w, newErrorString, http.StatusBadRequest)
		return
	}

	// Service
	err := h.ComputerService.AddComputer(newComputerDTO.Num, newComputerDTO.Price)
	if err != nil {
		if errors.Is(err, models.ErrPriceConflict) {
			http.Error(w, models.ErrPriceConflict.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, models.ErrComputerNumConflict) {
			http.Error(w, models.ErrComputerNumConflict.Error(), http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteComputerHandler(w http.ResponseWriter, r *http.Request) {
	var deleteComputerDTO DeleteComputerDTO
	if err := json.NewDecoder(r.Body).Decode(&deleteComputerDTO); err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, "Error: Incorrect data string.", http.StatusBadRequest)
			return
		}

		http.Error(w, newErrorString, http.StatusBadRequest)
		return
	}

	// Service
	if err := h.ComputerService.DeleteComputer(deleteComputerDTO.ID); err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, "Error: Incorrect data string.", http.StatusBadRequest)
			return
		}

		http.Error(w, newErrorString, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ChangeComputerPrice(w http.ResponseWriter, r *http.Request) {
	var changePriceDTO ChangePriceDTO
	if err := json.NewDecoder(r.Body).Decode(&changePriceDTO); err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, "Error: incorrect data string.", http.StatusBadRequest)
			return
		}

		http.Error(w, newErrorString, http.StatusBadRequest)
	}

	if err := h.ComputerService.ChangePrice(changePriceDTO.Number, changePriceDTO.NewPrice); err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, "Error: Incorrect data string.", http.StatusBadRequest)
			return
		}

		http.Error(w, newErrorString, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}
