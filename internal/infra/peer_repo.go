package infra

import (
	"context"
	"database/sql"
	"time"
)

type Peer struct {
	ID        int
	Config    string
	CreatedAt time.Time
}

type PeerRepo struct {
	db *sql.DB
}

func NewPeerRepo(db *sql.DB) *PeerRepo {
	return &PeerRepo{db: db}
}

func (r *PeerRepo) Save(ctx context.Context, id int, config string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO peers (id, config) VALUES ($1, $2)`,
		id,
		config,
	)
	return err
}

func (r *PeerRepo) Count(ctx context.Context) (int, error) {
	var cnt int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM peers`).Scan(&cnt)
	return cnt, err
}
