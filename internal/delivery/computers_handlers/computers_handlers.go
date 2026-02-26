package computers_handlers

import (
	"encoding/json"
	"mvp/internal/models"
	"mvp/internal/service/computers_service"
	"net/http"
	"time"
)

func AddComputerHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := computers_service.AddComputer(newComputerDTO); err != nil {
		newError := error.Error(err)
		http.Error(w, newError, http.StatusInternalServerError) // Статус код пока заглушка. Заменить.
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteComputerHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := computers_service.DeleteComputer(deleteComputerDTO); err != nil {
		newError := error.Error(err)
		http.Error(w, newError, http.StatusInternalServerError) // Статус код заглушка
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
