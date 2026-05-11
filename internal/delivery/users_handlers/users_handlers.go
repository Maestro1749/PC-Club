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
		switch {
		case errors.Is(err, models.ErrInternalServer):
			writeError(w, err, http.StatusInternalServerError)
		default:
			writeError(w, err, http.StatusConflict)
		}
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
		switch {
		case errors.Is(err, models.ErrNotFound):
			writeError(w, err, http.StatusNotFound)
		case errors.Is(err, models.ErrWrongPassword):
			writeError(w, err, http.StatusUnauthorized)
		default:
			writeError(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeError(w http.ResponseWriter, err error, status int) {
	newError := models.ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	str, _ := newError.ToString()

	http.Error(w, str, status)
}
