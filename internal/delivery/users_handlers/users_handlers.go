package users_handlers

import (
	"encoding/json"
	"errors"
	"mvp/internal/models"
	"mvp/internal/service/users_service"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type UserHandler struct {
	service *users_service.UserServise
	logger  *zap.Logger
}

func NewUserHandler(service *users_service.UserServise, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUserDTO models.NewUserDTO
	if err := json.NewDecoder(r.Body).Decode(&newUserDTO); err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, "Error: failed to formulate an error.", http.StatusInternalServerError)
			return
		}

		http.Error(w, newErrorString, http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterUser(newUserDTO); err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, "Error: failed to fotmulate an error.", http.StatusBadRequest)
			return
		}

		if errors.Is(err, models.ErrInternalServer) {
			http.Error(w, newErrorString, http.StatusInternalServerError)
			return
		}

		http.Error(w, newErrorString, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDTO models.LoginUserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, "Error: Incorrect data struct", http.StatusBadRequest)
			return
		}

		http.Error(w, newErrorString, http.StatusBadRequest)
		return
	}

	user, err := h.service.LoginUser(userDTO)
	if err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, newErrorString, http.StatusNotFound)
			return
		}

		if errors.Is(err, models.ErrWrongPassword) {
			http.Error(w, newErrorString, http.StatusUnauthorized)
			return
		}

		http.Error(w, newErrorString, http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
