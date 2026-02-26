package server

import (
	"errors"
	"mvp/internal/delivery/computers_handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer() error {
	router := mux.NewRouter()

	router.Path("/computers").Methods("CREATE").HandlerFunc(computers_handlers.AddComputerHandler)
	router.Path("/computers").Methods("DELETE").HandlerFunc(computers_handlers.DeleteComputerHandler)

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}
