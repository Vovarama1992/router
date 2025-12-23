package domain

import (
	"context"

	"router/internal/config"
	"router/internal/infra"
	"router/internal/wg"
)

type Peer struct {
	Config string
}

type Service struct {
	cfg   *config.Config
	repo  *infra.PeerRepo
	wgApp *infra.WGApplier
}

func NewService(cfg *config.Config, repo *infra.PeerRepo, wgApp *infra.WGApplier) *Service {
	return &Service{cfg: cfg, repo: repo, wgApp: wgApp}
}

func (s *Service) CreatePeer(ctx context.Context) (*Peer, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	clientIndex := count + 2

	p, err := wg.CreatePeer(s.cfg, clientIndex)
	if err != nil {
		return nil, err
	}

	// ВАЖНО: применяем peer к WireGuard
	if err := s.wgApp.ApplyPeer(ctx, p.PublicKey, p.Address); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, p.PublicKey, p.Address); err != nil {
		return nil, err
	}

	return &Peer{Config: p.Config}, nil
}
