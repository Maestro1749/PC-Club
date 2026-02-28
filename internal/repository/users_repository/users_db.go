package users_repository

import (
	"context"
	"database/sql"
	"errors"
	"mvp/internal/models"
	"time"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByPhoneNumber(phoneNumber string) (*models.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, fullname, email, phoneNumber, passwd, birthday, balance, privilage)
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
		user.Privilege,
	).Scan(&user.ID, &user.Registered)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, fullname, email, phoneNumber, passwd, birthday, balance, registered, privilage
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
			return nil, errors.New("User undefined.")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetByPhoneNumber(phoneNumber string) (*models.User, error) {
	query := `
		SELECT id, username, fullname, email, phoneNumber, passwd, birthday, balance, registered, privilage
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
			return nil, nil // user не найден, возможно, будет лучше возвращать сообщение.
		}
		return nil, err
	}

	return &user, nil
}
