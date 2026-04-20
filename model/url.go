package model

import "time"

type URL struct {
	ID          string
	OriginalURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   *time.Time
}
