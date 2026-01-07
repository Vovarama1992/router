package domain

import (
	"context"
	"fmt"
	"log"

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
	log.Println("[vpn] create peer: start")

	count, err := s.repo.Count(ctx)
	if err != nil {
		log.Printf("[vpn] count error: %v", err)
		return nil, err
	}

	name := fmt.Sprintf("peer_%d", count+1)
	log.Printf("[vpn] peer name: %s", name)

	// UDP-only — фиксируем явно
	client, err := openvpn.CreatePeer(name, "udp")
	if err != nil {
		log.Printf("[vpn] create peer failed: %v", err)
		return nil, err
	}

	if err := s.repo.Save(ctx, count+1, name); err != nil {
		log.Printf("[vpn] db save failed: %v", err)
		return nil, err
	}

	log.Printf("[vpn] peer created successfully: %s", name)

	return &Peer{
		Config: client.Config,
	}, nil
}
