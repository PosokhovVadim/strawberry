package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uint
	UserID    uint
	TokenID   uuid.UUID
	TokenHash string
	UserAgent string
	Device    string
	IP        string
	Location  string
	IsRevoked bool
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshSession struct {
	UserID    uint
	Token     string
	UserAgent string
	Device    string
	IP        string
	Location  string
}
