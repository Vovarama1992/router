package domain

import (
	"context"
	"fmt"
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

func (s *Service) CreatePeer(ctx context.Context) (*Peer, error) {
	start := time.Now()

	log.Printf("[domain] CreatePeer start")

	count, err := s.repo.Count(ctx)
	if err != nil {
		log.Printf("[domain] repo.Count FAILED err=%v", err)
		return nil, err
	}
	log.Printf("[domain] repo.Count OK count=%d", count)

	name := fmt.Sprintf("peer_%d", count+1)
	log.Printf("[domain] peer name=%s", name)

	client, err := reality.CreateClient()
	if err != nil {
		log.Printf("[domain] reality.CreateClient FAILED err=%v", err)
		return nil, err
	}

	log.Printf("[domain] reality.CreateClient OK uuid=%s", client.UUID)

	if err := s.repo.Save(ctx, count+1, client.UUID); err != nil {
		log.Printf("[domain] repo.Save FAILED id=%d uuid=%s err=%v", count+1, client.UUID, err)
		return nil, err
	}
	log.Printf("[domain] repo.Save OK id=%d uuid=%s", count+1, client.UUID)

	log.Printf("[domain] CreatePeer done uuid=%s duration=%s", client.UUID, time.Since(start))

	return &Peer{
		Link: client.Link,
	}, nil
}
