package models

import "time"

type User struct {
	ID          int
	TelegramID  int64
	ActiveUntil time.Time
	CreatedAt   time.Time
}
