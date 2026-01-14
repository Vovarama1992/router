package infra

import (
	"context"
	"database/sql"
	"log"
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
	start := time.Now()

	log.Printf("[repo] Save start id=%d", id)

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO peers (id, config) VALUES ($1, $2)`,
		id,
		config,
	)
	if err != nil {
		log.Printf("[repo] Save FAILED id=%d err=%v", id, err)
		return err
	}

	log.Printf("[repo] Save OK id=%d duration=%s", id, time.Since(start))
	return nil
}

func (r *PeerRepo) Count(ctx context.Context) (int, error) {
	start := time.Now()

	log.Printf("[repo] Count start")

	var cnt int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM peers`).Scan(&cnt)
	if err != nil {
		log.Printf("[repo] Count FAILED err=%v", err)
		return 0, err
	}

	log.Printf("[repo] Count OK count=%d duration=%s", cnt, time.Since(start))
	return cnt, nil
}
