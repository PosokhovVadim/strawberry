package model

import (
	"time"

	"github.com/PosokhovVadim/stawberry/internal/domain/entity"
	"github.com/google/uuid"
)

type Session struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index"`
	TokenID      uuid.UUID `gorm:"not null;unique;index"`
	Token        string    `gorm:"not null;unique;index"`
	Fingerprint  string
	UserAgent    string
	Device       string
	IP           string
	Location     string
	IsRevoked    bool
	RevokedAt    time.Time `gorm:"default:now();->"`
	ExpiresAt    time.Time
	LastActivity time.Time `gorm:"default:now()"`
	CreatedAt    time.Time `gorm:"default:now();->"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func ConvertSessionFromEntity(e entity.Session) Session {
	return Session{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenID:   e.TokenID,
		Token:     e.TokenHash,
		UserAgent: e.UserAgent,
		Device:    e.Device,
		IP:        e.IP,
		Location:  e.Location,
		IsRevoked: e.IsRevoked,
		ExpiresAt: e.ExpiresAt,
	}
}

func ConvertSessionToEntity(s Session) entity.Session {
	return entity.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		TokenID:   s.TokenID,
		TokenHash: s.Token,
		UserAgent: s.UserAgent,
		Device:    s.Device,
		IP:        s.IP,
		Location:  s.Location,
		IsRevoked: s.IsRevoked,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
