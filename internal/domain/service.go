package domain

import (
	"context"
	"fmt"
	"log"
	"time"

	"router/internal/infra"
	"router/internal/openvpn"
)

type Peer struct {
	Config string
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

	client, err := openvpn.CreatePeer(name, "tcp")
	if err != nil {
		log.Printf("[domain] openvpn.CreatePeer FAILED name=%s err=%v", name, err)
		return nil, err
	}
	log.Printf("[domain] openvpn.CreatePeer OK name=%s", name)

	if err := s.repo.Save(ctx, count+1, name); err != nil {
		log.Printf("[domain] repo.Save FAILED id=%d name=%s err=%v", count+1, name, err)
		return nil, err
	}
	log.Printf("[domain] repo.Save OK id=%d name=%s", count+1, name)

	log.Printf("[domain] CreatePeer done name=%s duration=%s", name, time.Since(start))

	return &Peer{
		Config: client.Config,
	}, nil
}
