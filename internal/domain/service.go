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

func (s *Service) CreatePeer(ctx context.Context) (*Peer, error) {
	start := time.Now()

	log.Printf("[domain] CreatePeer start")

	lastUUID, err := s.repo.GetLastUUID(ctx)
	if err != nil {
		log.Printf("[domain] repo.GetLastUUID FAILED err=%v", err)
		return nil, err
	}

	if lastUUID != "" {
		log.Printf("[domain] reuse existing uuid=%s", lastUUID)

		link, err := reality.BuildLink(lastUUID)
		if err != nil {
			log.Printf("[domain] BuildLink FAILED err=%v", err)
			return nil, err
		}

		return &Peer{
			Link: link,
		}, nil
	}

	client, err := reality.CreateClient()
	if err != nil {
		log.Printf("[domain] reality.CreateClient FAILED err=%v", err)
		return nil, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		log.Printf("[domain] repo.Count FAILED err=%v", err)
		return nil, err
	}

	if err := s.repo.Save(ctx, count+1, client.UUID); err != nil {
		log.Printf("[domain] repo.Save FAILED err=%v", err)
		return nil, err
	}

	log.Printf("[domain] CreatePeer done uuid=%s duration=%s", client.UUID, time.Since(start))

	return &Peer{
		Link: client.Link,
	}, nil
}
