package infra

import (
	"context"
	"database/sql"
	"time"

	"router/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByTelegramID(ctx context.Context, tgID int64) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, telegram_id, active_until, created_at
		FROM users
		WHERE telegram_id = $1
		LIMIT 1
	`, tgID)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.TelegramID,
		&u.ActiveUntil,
		&u.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, tgID int64) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (telegram_id)
		VALUES ($1)
	`, tgID)

	return err
}

func (r *UserRepo) IsActive(ctx context.Context, tgID int64) (bool, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT active_until
		FROM users
		WHERE telegram_id = $1
	`, tgID)

	var activeUntil time.Time

	err := row.Scan(&activeUntil)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return activeUntil.After(time.Now()), nil
}

func (r *UserRepo) UpdateActiveUntil(ctx context.Context, tgID int64, until time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET active_until = $2
		WHERE telegram_id = $1
	`, tgID, until)

	return err
}
