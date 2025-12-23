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
	cfg  *config.Config
	repo *infra.PeerRepo
}

func NewService(cfg *config.Config, repo *infra.PeerRepo) *Service {
	return &Service{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *Service) CreatePeer(ctx context.Context) (*Peer, error) {
	// 1) сколько пиров уже занято
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 2) следующий индекс (10.0.0.2, .3, ...)
	clientIndex := count + 2

	// 3) создаём peer
	p, err := wg.CreatePeer(s.cfg, clientIndex)
	if err != nil {
		return nil, err
	}

	// 4) сохраняем peer как занятый
	if err := s.repo.Save(ctx, p.PublicKey, p.Address); err != nil {
		return nil, err
	}

	// 5) отдаём конфиг клиенту
	return &Peer{
		Config: p.Config,
	}, nil
}
