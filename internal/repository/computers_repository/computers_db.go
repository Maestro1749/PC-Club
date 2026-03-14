package computers_repository

import (
	"context"
	"database/sql"
	"errors"
	"mvp/internal/models"
	"time"

	"github.com/lib/pq"
)

type ComputerRepository interface {
	CreateComputer(computer *models.Computer) error
	DeleteComputer(computer *models.Computer) error
	ChangePrice(num string, price float64) error
	GetByNumber(number string) (*models.Computer, error)
}

type computerRepo struct {
	db *sql.DB
}

func NewComputerRepository(db *sql.DB) ComputerRepository {
	return &computerRepo{db: db}
}

func (r *computerRepo) CreateComputer(computer *models.Computer) error {
	query := `
		INSERT INTO computers (num, price)
		VALUES ($1, $2)	
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		computer.Num,
		computer.Price,
	).Scan(&computer.ID)
	if err != nil {
		var pgErr *pq.Error

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				switch pgErr.Constraint {
				case "comuters_num_key":
					return models.ErrComputerNumConflict
				}
			}
		}

		return models.ErrInternalServer
	}

	return nil
}

func (r *computerRepo) DeleteComputer(computer *models.Computer) error {
	query := `
		DELETE FROM computers WHERE id = $1
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		computer.ID,
	).Scan(&computer.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *computerRepo) ChangePrice(num string, price float64) error {
	query := `
		UPDATE computers SET price = $1 WHERE num = $2
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var computer models.Computer

	err := r.db.QueryRowContext(
		ctx,
		query,
		price,
		num,
	).Scan(&computer.ID)
	if err != nil {
		return models.ErrInternalServer
	}

	return nil
}

func (r *computerRepo) GetByNumber(number string) (*models.Computer, error) {
	query := `
		SELECT id, num, price, is_busy, busy_since, busy_until
		FROM computers
		WHERE num = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var computer models.Computer

	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&computer.ID,
		&computer.Num,
		&computer.Price,
		&computer.IsBusy,
		&computer.BusySince,
		&computer.BusyUntil,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrComputerNotFound
		}
		return nil, err
	}

	return &computer, nil
}
