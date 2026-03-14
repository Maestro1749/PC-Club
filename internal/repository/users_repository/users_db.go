package users_repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"mvp/internal/models"
	"time"

	"github.com/lib/pq"
)

type UserRepository interface {
	Create(user *models.User) error
	LoginUser(username, password string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByPhoneNumber(phoneNumber string) (*models.User, error)
	GetHashByUsername(username string) ([]byte, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, fullname, email, phone_number, passwd, birthday, balance, registered)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, registered
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Fullname,
		user.Email,
		user.PhoneNumber,
		user.Password,
		user.Birthday,
		user.Balance,
		user.Registered,
	).Scan(&user.ID, &user.Registered)
	if err != nil {
		log.Printf("DB ERROR: %v\n", err)

		var pgErr *pq.Error

		if errors.As(err, &pgErr) {
			log.Printf("PG ERROR: %v %v\n", pgErr.Code, pgErr.Constraint)

			if pgErr.Code == "23505" {
				switch pgErr.Constraint {
				case "users_username_key":
					return models.ErrUsernameConflict
				case "users_phone_number_key":
					return models.ErrPhoneNumberConflict
				case "users_email_key":
					return models.ErrEmailConflict
				}
			}
		}

		return models.ErrInternalServer
	}

	return nil
}

func (r *userRepo) LoginUser(login string, passwordHash string) (*models.User, error) {
	query := `
		SELECT id, username, fullname, email, phone_number, birthday, balance, registered, privilage
		FROM users
		WHERE username = $1 AND passwd = $2
	`

	user := models.User{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		login,
		passwordHash,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Fullname,
		&user.Email,
		&user.PhoneNumber,
		&user.Birthday,
		&user.Balance,
		&user.Registered,
		&user.Privilege,
	)
	if err != nil {
		return nil, models.ErrInternalServer // Ошибка на стороне сервера
	}

	return &user, nil
}

func (r *userRepo) GetHashByUsername(username string) ([]byte, error) {
	query := `
		SELECT passwd FROM users
		WHERE username = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var hash []byte

	err := r.db.QueryRowContext(ctx, query, username).Scan(&hash)
	if err != nil {
		return nil, models.ErrNotFound
	}

	return hash, nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, fullname, email, phone_number, passwd, birthday, balance, registered, privilage
		FROM users
		WHERE email = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Fullname,
		&user.Email,
		&user.PhoneNumber,
		&user.Password,
		&user.Birthday,
		&user.Balance,
		&user.Registered,
		&user.Privilege,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrInternalServer
	}

	return &user, nil
}

func (r *userRepo) GetByPhoneNumber(phoneNumber string) (*models.User, error) {
	query := `
		SELECT id, username, fullname, email, phone_number, passwd, birthday, balance, registered, privilage
		FROM users
		WHERE phone_number = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.ID,
		&user.Username,
		&user.Fullname,
		&user.Email,
		&user.PhoneNumber,
		&user.Password,
		&user.Birthday,
		&user.Balance,
		&user.Registered,
		&user.Privilege,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrInternalServer
	}

	return &user, nil
}
