package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"uid"`
	Role   uint `json:"rol"`
	Type   uint `json:"typ"`

	jwt.RegisteredClaims
}
