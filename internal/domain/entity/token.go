package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uint `json:"uid"`
	Role   uint `json:"rol"`
	Type   uint `json:"typ"`

	jwt.RegisteredClaims
}

type Token struct {
	JTI       uuid.UUID `json:"-"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TokenPair struct {
	Access  Token `json:"access"`
	Refresh Token `json:"refresh"`
}

type Refresh struct {
	Token string `json:"token"`
}

type AccessTokens struct {
	Access           string    `json:"access"`
	Refresh          string    `json:"refresh"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}
