package domain

import (
	"context"
	"log"
	"time"

	"router/internal/infra"
	"router/internal/reality"
)

type Peer struct {
	Link string
}

type Service struct {
	repo *infra.PeerRepo
}

func NewService(repo *infra.PeerRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePeer(ctx context.Context, telegramID int64) (*Peer, error) {
	start := time.Now()

	log.Printf("[domain] CreatePeer start tg=%d", telegramID)

	// 1 — проверяем есть ли уже
	existing, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		link, err := reality.BuildLink(existing.UUID)
		if err != nil {
			return nil, err
		}

		log.Printf("[domain] reuse uuid=%s", existing.UUID)

		return &Peer{Link: link}, nil
	}

	// 2 — создаём нового клиента
	client, err := reality.CreateClient()
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, client.UUID, telegramID); err != nil {
		return nil, err
	}

	log.Printf("[domain] CreatePeer done uuid=%s duration=%s", client.UUID, time.Since(start))

	return &Peer{
		Link: client.Link,
	}, nil
}

type PeerInfo struct {
	TelegramID int64  `json:"telegram_id"`
	UUID       string `json:"uuid"`
	Status     string `json:"status"`
}

func (s *Service) ListPeers(ctx context.Context) ([]PeerInfo, error) {
	peers, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var out []PeerInfo
	for _, p := range peers {
		out = append(out, PeerInfo{
			TelegramID: p.TelegramID,
			UUID:       p.UUID,
			Status:     p.ConnectionStatus,
		})
	}

	return out, nil
}

func (s *Service) DisableByTelegramID(ctx context.Context, telegramID int64) error {
	// 1 — получаем все UUID пользователя
	peers, err := s.repo.ListByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	if len(peers) == 0 {
		return nil
	}

	var uuids []string
	for _, p := range peers {
		uuids = append(uuids, p.UUID)
	}

	// 2 — удалить из config
	if err := reality.RemoveClients(uuids); err != nil {
		return err
	}

	return nil
}
