package domain

import (
	"context"

	"router/internal/config"
	"router/internal/infra"
	"router/internal/openvpn"
)

type Peer struct {
	Config string
}

type Service struct {
	cfg  *config.Config
	repo *infra.PeerRepo
}

func NewService(cfg *config.Config, repo *infra.PeerRepo) *Service {
	return &Service{cfg: cfg, repo: repo}
}

func (s *Service) CreatePeer(ctx context.Context) (*Peer, error) {
	// 1. Считаем текущее количество
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Генерим клиентский конфиг
	client, err := openvpn.CreateClient()
	if err != nil {
		return nil, err
	}

	// 3. Сохраняем в БД (номер + конфиг)
	if err := s.repo.Save(ctx, count+1, client.Config); err != nil {
		return nil, err
	}

	// 4. Отдаем клиенту
	return &Peer{
		Config: client.Config,
	}, nil
}
