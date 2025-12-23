package infra

import (
	"context"
	"database/sql"
)

type Peer struct {
	ID         int
	PublicKey  string
	VPNAddress string
}

type PeerRepo struct {
	db *sql.DB
}

func NewPeerRepo(db *sql.DB) *PeerRepo {
	return &PeerRepo{db: db}
}

func (r *PeerRepo) Save(ctx context.Context, publicKey, vpnAddress string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO peers (public_key, vpn_address) VALUES ($1, $2)`,
		publicKey,
		vpnAddress,
	)
	return err
}

func (r *PeerRepo) Count(ctx context.Context) (int, error) {
	var cnt int
	err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM peers`).Scan(&cnt)
	return cnt, err
}
