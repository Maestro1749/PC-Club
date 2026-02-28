package computers_handlers

import (
	"encoding/json"
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
	var newComputerDTO models.NewComputerDTO
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

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteComputerHandler(w http.ResponseWriter, r *http.Request) {
	var deleteComputerDTO models.DeleteComputerDTO
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

	w.WriteHeader(http.StatusNoContent)
}
