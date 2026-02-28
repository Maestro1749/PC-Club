package main

import (
	"database/sql"
	"log"
	"mvp/internal/delivery/computers_handlers"
	"mvp/internal/repository/computers_repository"
	"mvp/internal/service/computers_service"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Repositories
	//userRepo := users_repository.NewUserRepository(db)
	computerRepo := computers_repository.NewComputerRepository(db)

	// Services
	computerService := computers_service.NewComputerService(computerRepo)

	// Handlers
	computerHandler := computers_handlers.NewHandler(computerService)

	router := mux.NewRouter()

	router.Path("/register").Methods("CREATE").HandlerFunc(computerHandler.AddComputerHandler)

	if err := http.ListenAndServe(":9091", router); err != nil {
		log.Fatal(err)
	}
}
