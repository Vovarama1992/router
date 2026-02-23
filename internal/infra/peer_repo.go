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

func (r *PeerRepo) GetLastUUID(ctx context.Context) (string, error) {
	start := time.Now()

	log.Printf("[repo] GetLastUUID start")

	var uuid string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT config FROM peers ORDER BY id DESC LIMIT 1`,
	).Scan(&uuid)

	if err == sql.ErrNoRows {
		log.Printf("[repo] GetLastUUID empty duration=%s", time.Since(start))
		return "", nil
	}

	if err != nil {
		log.Printf("[repo] GetLastUUID FAILED err=%v", err)
		return "", err
	}

	log.Printf("[repo] GetLastUUID OK uuid=%s duration=%s", uuid, time.Since(start))

	return uuid, nil
}

func (r *PeerRepo) GetByID(ctx context.Context, id int) (string, error) {
	start := time.Now()

	log.Printf("[repo] GetByID start id=%d", id)

	var config string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT config FROM peers WHERE id = $1`,
		id,
	).Scan(&config)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[repo] GetByID no rows id=%d", id)
			return "", nil
		}
		log.Printf("[repo] GetByID FAILED id=%d err=%v", id, err)
		return "", err
	}

	log.Printf("[repo] GetByID OK id=%d duration=%s", id, time.Since(start))
	return config, nil
}
