package usershandlers

import (
	"encoding/json"
	"mvp/internal/models"
	"net/http"
	"time"
)

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
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
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

	// Вызов сервиса

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
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusOK)
}
