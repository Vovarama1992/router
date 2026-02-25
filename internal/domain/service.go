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
