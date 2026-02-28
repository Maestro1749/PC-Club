package main

import (
	"database/sql"
	"fmt"
	"log"
	"mvp/internal/delivery/computers_handlers"
	"mvp/internal/repository/computers_repository"
	"mvp/internal/service/computers_service"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connectString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connectString)
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

	router.Path("/add").Methods("POST").HandlerFunc(computerHandler.AddComputerHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
