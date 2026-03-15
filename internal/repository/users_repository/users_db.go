package users_repository

import (
	"context"
	"database/sql"
	"errors"
	"mvp/internal/models"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type UserRepository interface {
	Create(user *models.User) error
	LoginUser(username, password string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByPhoneNumber(phoneNumber string) (*models.User, error)
	GetHashByUsername(username string) ([]byte, error)
}

type userRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) UserRepository {
	return &userRepo{
		db:     db,
		logger: logger,
	}
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
		var pgErr *pq.Error

		if errors.As(err, &pgErr) {
			r.logger.Error("Failed to create user", zap.Error(err))

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

		r.logger.Error("Failed to create user", zap.Error(err))
		return models.ErrInternalServer
	}

	r.logger.Info("User succesfully registered", zap.String("username", user.Username))
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		r.logger.Error("Failed to execute login query", zap.Error(err))
		return nil, models.ErrInternalServer
	}

	r.logger.Info("User successfully login", zap.String("username", user.Username))
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		r.logger.Error("Failed to execute get hash query", zap.Error(err))
		return nil, models.ErrInternalServer
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
		r.logger.Error("Failed to execute get by email query", zap.Error(err))
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
		r.logger.Error("Failed to execute get by phone number query", zap.Error(err))
		return nil, models.ErrInternalServer
	}

	return &user, nil
}
