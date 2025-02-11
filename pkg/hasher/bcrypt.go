package hasher

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type BCryptHasher struct {
	cfg *Config
}

func NewBcryptHasher(cfg *Config) (*BCryptHasher, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &BCryptHasher{cfg: cfg}, nil
}

func (h *BCryptHasher) Hash(_ context.Context, value string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), int(h.cfg.Cost))
	return string(bytes), err
}

func (h *BCryptHasher) Compare(_ context.Context, hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
