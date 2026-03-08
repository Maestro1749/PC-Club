package users_handlers

import (
	"encoding/json"
	"mvp/internal/models"
	"mvp/internal/service/users_service"
	"net/http"
	"time"
)

type UserHandler struct {
	service *users_service.UserServise
}

func NewUserHandler(service *users_service.UserServise) *UserHandler {
	return &UserHandler{service: service}
}

/*
pattern: /users
method: CREATE
info: JSON in HTTP request body

succeed:
  - status code: 201 Created
  - response body: JSON represent created user

failed:
  - status code: 400, 409, 500, ...
  - response body: JSON represent user
*/
func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUserDTO models.NewUserDTO
	if err := json.NewDecoder(r.Body).Decode(&newUserDTO); err != nil {
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

	if err := h.service.RegisterUser(newUserDTO); err != nil {
		newError := models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		newErrorString, err := newError.ToString()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, newErrorString, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

/*
pattern: /users
method: GET
info: JSON in HTTP request body

succeed:
  - status code: 200 OK
  - responce body: JSON represent user

failed:
  - status code: 400, 409, ...
  - response body: JSON error
*/
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

		http.Error(w, newErrorString, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
