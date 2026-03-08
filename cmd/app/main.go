package main

import (
	"database/sql"
	"fmt"
	"log"
	"mvp/internal/delivery/computers_handlers"
	"mvp/internal/delivery/users_handlers"
	"mvp/internal/repository/computers_repository"
	"mvp/internal/repository/users_repository"
	"mvp/internal/service/computers_service"
	"mvp/internal/service/users_service"
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
	userRepo := users_repository.NewUserRepository(db)

	// Services
	computerService := computers_service.NewComputerService(computerRepo)
	userService := users_service.NewService(userRepo)

	// Handlers
	computerHandler := computers_handlers.NewHandler(computerService)
	userHandler := users_handlers.NewUserHandler(userService)

	router := mux.NewRouter()

	router.Path("/computer/add").Methods("POST").HandlerFunc(computerHandler.AddComputerHandler)
	router.Path("/computer/delete").Methods("DELETE").HandlerFunc(computerHandler.DeleteComputerHandler)
	router.Path("/computer/changeprice").Methods("PUT").HandlerFunc(computerHandler.ChangeComputerPrice)

	router.Path("/user/register").Methods("POST").HandlerFunc(userHandler.RegisterUserHandler)
	router.Path("/user/login").Methods("GET").HandlerFunc(userHandler.LoginUserHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
