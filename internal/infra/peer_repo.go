package infra

import (
	"context"
	"database/sql"
	"router/internal/models"
)

type PeerRepo struct {
	db *sql.DB
}

func NewPeerRepo(db *sql.DB) *PeerRepo {
	return &PeerRepo{db: db}
}

func (r *PeerRepo) GetByTelegramID(ctx context.Context, tgID int64) (*models.Peer, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, uuid, telegram_id, is_active, created_at
		FROM peers
		WHERE telegram_id = $1
		LIMIT 1
	`, tgID)

	var p models.Peer
	err := row.Scan(&p.ID, &p.UUID, &p.TelegramID, &p.IsActive, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PeerRepo) Create(ctx context.Context, uuid string, tgID int64) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO peers (uuid, telegram_id, is_active)
		VALUES ($1, $2, TRUE)
	`, uuid, tgID)
	return err
}

func (r *PeerRepo) List(ctx context.Context) ([]models.Peer, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, uuid, telegram_id, is_active, created_at
		FROM peers
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []models.Peer

	for rows.Next() {
		var p models.Peer
		if err := rows.Scan(
			&p.ID,
			&p.UUID,
			&p.TelegramID,
			&p.IsActive,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		peers = append(peers, p)
	}

	return peers, nil
}

func (r *PeerRepo) ListByTelegramID(ctx context.Context, tgID int64) ([]models.Peer, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, uuid, telegram_id, is_active, created_at
		FROM peers
		WHERE telegram_id = $1
	`, tgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []models.Peer

	for rows.Next() {
		var p models.Peer
		if err := rows.Scan(
			&p.ID,
			&p.UUID,
			&p.TelegramID,
			&p.IsActive,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		peers = append(peers, p)
	}

	return peers, nil
}

func (r *PeerRepo) SetActive(ctx context.Context, tgID int64, active bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE peers
		SET is_active = $1
		WHERE telegram_id = $2
	`, active, tgID)
	return err
}

func (r *PeerRepo) Reactivate(ctx context.Context, tgID int64) error {
	return r.SetActive(ctx, tgID, true)
}
