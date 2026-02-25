package models

import "time"

type Peer struct {
	ID               int
	UUID             string
	TelegramID       int64
	ConnectionStatus string
	CreatedAt        time.Time
}
