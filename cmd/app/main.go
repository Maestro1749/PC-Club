package main

import (
	"database/sql"
	"fmt"
	"mvp/internal/delivery/booking_handlers"
	"mvp/internal/delivery/computers_handlers"
	"mvp/internal/delivery/users_handlers"
	"mvp/internal/logger"
	"mvp/internal/repository/booking_repository"
	"mvp/internal/repository/computers_repository"
	"mvp/internal/repository/users_repository"
	"mvp/internal/service/booking_service"
	"mvp/internal/service/computers_service"
	"mvp/internal/service/users_service"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	logger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("Logger initialized successfully")

	connectString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		logger.Error("Failed to open database connection", zap.Error(err))
		panic(err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("Failed to connect to the database", zap.Error(err))
		panic(err)
	}
	logger.Info("Database connection established successfully")

	// Repositories
	computerRepo := computers_repository.NewComputerRepository(db, logger)
	userRepo := users_repository.NewUserRepository(db, logger)
	bookingRepo := booking_repository.NewBookingReposiry(db, logger)
	logger.Info("Repositories initialized successfully")

	// Services
	computerService := computers_service.NewComputerService(computerRepo, logger)
	userService := users_service.NewService(userRepo, logger)
	bookingService := booking_service.NewBookingService(bookingRepo, logger)
	logger.Info("Services initialized successfully")

	// Handlers
	computerHandler := computers_handlers.NewHandler(computerService, logger)
	userHandler := users_handlers.NewUserHandler(userService, logger)
	bookingHandler := booking_handlers.NewBookingHandler(bookingService, logger)
	logger.Info("Handlers initialized successfully")

	router := mux.NewRouter()

	router.Path("/computer/add").Methods("POST").HandlerFunc(computerHandler.AddComputerHandler)
	router.Path("/computer/delete").Methods("DELETE").HandlerFunc(computerHandler.DeleteComputerHandler)
	router.Path("/computer/changeprice").Methods("PUT").HandlerFunc(computerHandler.ChangeComputerPrice)

	router.Path("/user/register").Methods("POST").HandlerFunc(userHandler.RegisterUserHandler)
	router.Path("/user/login").Methods("GET").HandlerFunc(userHandler.LoginUserHandler)

	router.Path("/booking").Methods("POST").HandlerFunc(bookingHandler.BookingComputerHandler)

	logger.Info("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("Failed to start server", zap.Error(err))
		panic(err)
	}
}
