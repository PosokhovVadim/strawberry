package session

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/PosokhovVadim/stawberry/internal/app/apperror"
	"github.com/PosokhovVadim/stawberry/internal/domain/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	SessionRepository interface {
		Create(ctx context.Context, session entity.Session) (sessionID uint, err error)
		Get(ctx context.Context, tokenID uuid.UUID) (session entity.Session, err error)
		GetForUpdate(ctx context.Context, tokenID uuid.UUID) (session entity.Session, err error)
		Update(ctx context.Context, session entity.Session) (err error)

		Transaction(ctx context.Context, fn func(context.Context) error) error
	}

	Dependencies struct {
		Config            *Config
		SessionRepository SessionRepository
	}

	TokenConfig struct {
		Type   TokenType
		TTL    time.Duration
		Secret *ecdsa.PrivateKey
	}

	Config struct {
		Access  TokenConfig
		Refresh TokenConfig
	}

	Service struct {
		config      *Config
		sessionRepo SessionRepository
	}
)

type TokenType uint

const (
	TokenAccessType TokenType = iota
	TokenRefreshType
	TokenUnknown
)

func generateESDSAKey() *ecdsa.PrivateKey {
	secret, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(fmt.Errorf("failed to generate secret: %w", err))
	}

	return secret
}

var DefaultConfig = &Config{
	Access: TokenConfig{
		Type:   TokenAccessType,
		TTL:    15 * time.Minute,
		Secret: generateESDSAKey(),
	},
	Refresh: TokenConfig{
		Type:   TokenRefreshType,
		TTL:    7 * 24 * time.Hour,
		Secret: generateESDSAKey(),
	},
}

func New(deps *Dependencies) *Service {
	s := &Service{
		config:      deps.Config,
		sessionRepo: deps.SessionRepository,
	}

	if s.config == nil {
		s.config = DefaultConfig
	}

	return s
}

func (s *Service) Create(ctx context.Context, session entity.Session) (tokens entity.TokenPair, err error) {
	if tokens, err = s.prepareTokens(ctx, session.UserID); err != nil {
		return entity.TokenPair{}, apperror.ErrSessionInternalError.Internal(err)
	}

	newSession := entity.Session{
		UserID:    session.UserID,
		UserAgent: session.UserAgent,
		Device:    session.Device,
		TokenID:   tokens.Refresh.JTI,
		IP:        session.IP,
		Location:  session.Location,
		TokenHash: tokens.Refresh.Token,
		ExpiresAt: tokens.Refresh.ExpiresAt,
	}

	if _, err = s.sessionRepo.Create(ctx, newSession); err != nil {
		return entity.TokenPair{}, apperror.ErrSessionDatabaseError.Internal(err)
	}

	return tokens, nil
}

func (s *Service) Refresh(ctx context.Context, req entity.RefreshSession) (tokens entity.TokenPair, err error) {
	var claims entity.Claims
	if claims, err = s.verify(ctx, req.Token, s.config.Refresh); err != nil {
		return entity.TokenPair{}, err
	}

	err = s.sessionRepo.Transaction(ctx, func(ctx context.Context) error {
		var activeSession entity.Session
		if activeSession, err = s.sessionRepo.GetForUpdate(ctx, uuid.MustParse(claims.ID)); err != nil {
			if errors.Is(err, apperror.ErrSessionNotFound) {
				return apperror.ErrSessionInvalidToken
			}
			return apperror.ErrSessionDatabaseError.Internal(err)
		}

		if err = s.validateSession(ctx, activeSession, req); err != nil {
			return err
		}

		if tokens, err = s.prepareTokens(ctx, activeSession.UserID); err != nil {
			return err
		}

		activeSession.IsRevoked = true
		if err = s.sessionRepo.Update(ctx, activeSession); err != nil {
			return apperror.ErrSessionDatabaseError.Internal(err)
		}

		newSession := entity.Session{
			UserID:    activeSession.UserID,
			UserAgent: req.UserAgent,
			Device:    req.Device,
			IP:        req.IP,
			Location:  req.Location,
			TokenID:   tokens.Refresh.JTI,
			TokenHash: tokens.Refresh.Token,
			ExpiresAt: tokens.Refresh.ExpiresAt,
		}

		if _, err = s.sessionRepo.Create(ctx, newSession); err != nil {
			return apperror.ErrSessionDatabaseError.Internal(err)
		}

		return nil
	})
	if err != nil {
		return entity.TokenPair{}, apperror.ErrSessionDatabaseTransactionError.Internal(err)
	}

	return tokens, nil
}

func (s *Service) validateSession(_ context.Context, c entity.Session, r entity.RefreshSession) (err error) {
	if c.TokenHash != r.Token {
		return apperror.ErrSessionInvalidToken
	}

	if c.IsRevoked == true {
		// Revoke all user active sessions?
		return apperror.ErrSessionSecurityViolation.Internal(errors.New("reusing revoked token"))
	}

	if time.Now().After(c.ExpiresAt) {
		return apperror.ErrSessionExpired
	}

	// Check device fingerprint, user-agent, may be ip, location also..
	// apperror.ErrSessionSecurityViolation

	return nil
}

func (s *Service) prepareTokens(ctx context.Context, userID uint) (tokens entity.TokenPair, err error) {
	if tokens.Access, err = s.generateToken(ctx, s.config.Access, entity.Claims{UserID: userID}); err != nil {
		return entity.TokenPair{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	if tokens.Refresh, err = s.generateToken(ctx, s.config.Refresh, entity.Claims{UserID: userID}); err != nil {
		return entity.TokenPair{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	return tokens, nil
}

func (s *Service) generateToken(_ context.Context, cfg TokenConfig, claims entity.Claims) (token entity.Token, err error) {
	issuedAt := time.Now()
	token.ExpiresAt = issuedAt.Add(cfg.TTL)

	token.JTI, err = uuid.NewV7()
	if err != nil {
		return entity.Token{}, fmt.Errorf("failed to generate token ID: %w", err)
	}

	claims.Type = uint(cfg.Type)
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ID:        token.JTI.String(),
		Issuer:    "strawberry",
		ExpiresAt: jwt.NewNumericDate(token.ExpiresAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	if token.Token, err = jwtToken.SignedString(cfg.Secret); err != nil {
		return entity.Token{}, fmt.Errorf("failed to signed token: %w", err)
	}

	return token, nil
}

func (s *Service) VerifyAccessToken(ctx context.Context, token string) (claims entity.Claims, err error) {
	return s.verify(ctx, token, s.config.Access)
}

func (s *Service) VerifyRefresh(ctx context.Context, token string) (claims entity.Claims, err error) {
	return s.verify(ctx, token, s.config.Refresh)
}

func (s *Service) verify(_ context.Context, token string, cfg TokenConfig) (claims entity.Claims, err error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodES256.Name}),
		jwt.WithIssuer("strawberry"),
		jwt.WithExpirationRequired(),
	)

	verifiedClaims := entity.Claims{}
	jwtToken, err := parser.ParseWithClaims(token, &verifiedClaims, func(tok *jwt.Token) (any, error) {
		return cfg.Secret.Public(), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return entity.Claims{}, apperror.ErrSessionExpired
		default:
			return entity.Claims{}, apperror.ErrSessionInvalidToken
		}
	}

	if !jwtToken.Valid {
		return entity.Claims{}, fmt.Errorf("token invalid")
	}

	return verifiedClaims, nil
}
