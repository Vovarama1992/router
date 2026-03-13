package domain

import (
	"context"
	"errors"
	"time"

	"router/internal/infra"
	"router/internal/reality"
)

var ErrAccessDisabled = errors.New("access disabled")

type Peer struct {
	Link string
}

type PeerInfo struct {
	TelegramID int64  `json:"telegram_id"`
	UUID       string `json:"uuid"`
	IsActive   bool   `json:"is_active"`
}

type Service struct {
	peerRepo *infra.PeerRepo
	userRepo *infra.UserRepo
}

func NewService(peerRepo *infra.PeerRepo, userRepo *infra.UserRepo) *Service {
	return &Service{
		peerRepo: peerRepo,
		userRepo: userRepo,
	}
}

func (s *Service) CreatePeer(ctx context.Context, telegramID int64) (*Peer, error) {

	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if err := s.userRepo.Create(ctx, telegramID); err != nil {
			return nil, err
		}
	}

	active, err := s.userRepo.IsActive(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if !active {
		return nil, ErrAccessDisabled
	}

	existing, err := s.peerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		link, err := reality.BuildLink(existing.UUID)
		if err != nil {
			return nil, err
		}

		return &Peer{Link: link}, nil
	}

	client, err := reality.CreateClient()
	if err != nil {
		return nil, err
	}

	if err := s.peerRepo.Create(ctx, client.UUID, telegramID); err != nil {
		return nil, err
	}

	return &Peer{Link: client.Link}, nil
}

func (s *Service) ListPeers(ctx context.Context) ([]PeerInfo, error) {
	peers, err := s.peerRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	var out []PeerInfo
	for _, p := range peers {
		out = append(out, PeerInfo{
			TelegramID: p.TelegramID,
			UUID:       p.UUID,
			IsActive:   p.IsActive,
		})
	}

	return out, nil
}

func (s *Service) SetUserUntil(ctx context.Context, telegramID int64, until time.Time) error {

	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	if user == nil {
		if err := s.userRepo.Create(ctx, telegramID); err != nil {
			return err
		}
	}

	return s.userRepo.UpdateActiveUntil(ctx, telegramID, until)
}
