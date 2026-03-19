package booking_repository

import (
	"context"
	"database/sql"
	"errors"
	"mvp/internal/models"
	"time"

	"go.uber.org/zap"
)

type BookingRepository interface {
	Booking(userID int, sum float64, computerID int, timeStart time.Time, timeEnd time.Time) error
	Exists(number string) (bool, error)
	IsAvailable(computerID int, timeStart, timeEnd time.Time) (bool, error)
	FindComputerIDByNumber(number string) (int, error)
	BusyCheck(timeStart time.Time, timeEnd time.Time, computerID int) (bool, error)
	TakePrice(computerNumber string) (float64, error)
}

type bookingRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewBookingReposiry(db *sql.DB, logger *zap.Logger) BookingRepository {
	return &bookingRepo{
		db:     db,
		logger: logger,
	}
}

func (r *bookingRepo) Booking(userID int, sum float64, computerID int, timeStart time.Time, timeEnd time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to start transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	query_users := `
		UPDATE users SET balance = balance - $1 WHERE id = $2 AND balance >= $1;
	`
	res, err := tx.ExecContext(ctx, query_users, sum, userID)
	if err != nil {
		r.logger.Error("payment error", zap.Error(err))
		return models.ErrInternalServer
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		r.logger.Info("low balance or other error", zap.Int("user id:", userID))
		return models.ErrInternalServer
	}

	query_computers := `
		UPDATE computers SET busy_since = $1, busy_until = $2 WHERE id = $3;
	`
	_, err = tx.ExecContext(ctx, query_computers, timeStart, timeEnd, computerID)
	if err != nil {
		r.logger.Error("update computers table error", zap.Error(err))
		return models.ErrInternalServer
	}

	query_booking := `
		INSERT INTO booking (user_id, computer_id, time_start, time_end)
		VALUES ($1, $2, $3, $4);
	`
	_, err = tx.ExecContext(ctx, query_booking, userID, computerID, timeStart, timeEnd)
	if err != nil {
		r.logger.Error("failed to write booking info in database", zap.Error(err))
		return models.ErrInternalServer
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return models.ErrInternalServer
	}

	return nil
}

func (r *bookingRepo) Exists(number string) (bool, error) {
	query := `
		SELECT 1 FROM computers WHERE num = $1 LIMIT 1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exist int

	err := r.db.QueryRowContext(ctx, query, number).Scan(&exist)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		r.logger.Error("Failed to check computer existence", zap.String("computer number", number), zap.Error(err))
		return false, models.ErrInternalServer
	}
	return true, nil
}

func (r *bookingRepo) IsAvailable(computerID int, timeStart, timeEnd time.Time) (bool, error) {
	query := `
		SELECT 1 FROM bookings
		WHERE cmputer_id = $1
		AND (time_start < $3 AND time_end > $2)
		LIMIT 1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exist int
	err := r.db.QueryRowContext(ctx, query, computerID, timeStart, timeEnd).Scan(&exist)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		r.logger.Error("Failed to check computer availability", zap.Int("computerID", computerID), zap.Time("timeStart", timeStart), zap.Time("timeEnd", timeEnd), zap.Error(err))
		return false, models.ErrInternalServer
	}
	return false, nil
}

func (r *bookingRepo) FindComputerIDByNumber(number string) (int, error) {
	query := `
		SELECT id FROM computers WHERE num = $1 LIMIT 1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var computerID int
	err := r.db.QueryRowContext(ctx, query, number).Scan(&computerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrComputerNotFound
		}
		r.logger.Error("Failed to find computer by number", zap.String("computer number", number), zap.Error(err))
		return 0, models.ErrInternalServer
	}
	return computerID, nil
}

func (r *bookingRepo) BusyCheck(timeStart time.Time, timeEnd time.Time, computerID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT 1 FROM booking
			WHERE computer_id = $1
			AND (time_start < $3 AND time_end > $2)
		);
	`

	var isBusy bool

	if err := r.db.QueryRowContext(ctx, query, computerID, timeStart, timeEnd).Scan(&isBusy); err != nil {
		r.logger.Error("error to check booking", zap.Error(err))
		return true, models.ErrInternalServer
	}

	return isBusy, nil
}

func (r *bookingRepo) TakePrice(computerNumber string) (float64, error) {
	query := `
		SELECT price FROM computers WHERE num = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var price float64
	if err := r.db.QueryRowContext(ctx, query, computerNumber).Scan(&price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Info("computer not found")
			return 0, models.ErrNotFound
		}

		r.logger.Error("error to get price", zap.Error(err))
		return 0, models.ErrInternalServer
	}

	return price, nil
}
