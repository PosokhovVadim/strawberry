package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/PosokhovVadim/stawberry/internal/app/apperror"
	"github.com/PosokhovVadim/stawberry/internal/domain/entity"
)

type PasswordPolicy struct {
	MinLength    int
	RequireUpper bool
	RequireLower bool
	RequireDigit bool
}

var (
	DefaultPasswordPolicy = &PasswordPolicy{
		MinLength:    8,
		RequireUpper: true,
		RequireLower: true,
		RequireDigit: true,
	}

	DefaultConfig = &Config{
		CheckPassword:    true,
		CheckEmailFormat: true,
		CheckEmailMX:     false,
	}

	requireUpperRegexp = regexp.MustCompile(`[A-Z]`)
	requireLowerRegexp = regexp.MustCompile(`[a-z]`)
	requireDigitRegexp = regexp.MustCompile(`[0-9]`)

	// Reference: https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
	emailFormatRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$") //nolint:lll
)

type (
	UserRepository interface {
		GetByEmail(ctx context.Context, email string) (user entity.User, err error)
		Save(ctx context.Context, register entity.Register) (userID uint, err error)

		Transaction(ctx context.Context, fn func(context.Context) error) error
	}

	SessionManager interface {
		Create(ctx context.Context, session entity.Session) (tokens entity.TokenPair, err error)
		Refresh(ctx context.Context, refresh entity.RefreshSession) (tokens entity.TokenPair, err error)
	}

	PasswordHasher interface {
		Hash(ctx context.Context, password string) (hash string, err error)
		Compare(ctx context.Context, hash string, password string) (err error)
	}

	EmailVerifier interface {
		VerifyMX(ctx context.Context, email string) (err error)
	}
)

type (
	Config struct {
		CheckPassword    bool
		CheckEmailFormat bool
		CheckEmailMX     bool
	}

	Dependencies struct {
		Config         *Config
		UserRepo       UserRepository
		SessionManager SessionManager
		PasswordHasher PasswordHasher
		PasswordPolicy *PasswordPolicy
		DNSResolver    EmailVerifier
	}

	Service struct {
		config         *Config
		userRepo       UserRepository
		sessionManager SessionManager
		passwordHasher PasswordHasher
		passwordPolicy *PasswordPolicy
		emailVerifier  EmailVerifier
	}
)

func New(deps *Dependencies) *Service {
	svc := &Service{
		config:         deps.Config,
		userRepo:       deps.UserRepo,
		sessionManager: deps.SessionManager,
		passwordHasher: deps.PasswordHasher,
		emailVerifier:  deps.DNSResolver,
	}

	if svc.config == nil {
		svc.config = DefaultConfig
	}

	if svc.passwordPolicy == nil {
		svc.passwordPolicy = DefaultPasswordPolicy
	}

	return svc
}

func (s *Service) Register(ctx context.Context, register entity.Register) (tokens entity.TokenPair, err error) {
	// Fail-fast with check email already exists..
	if err = s.emailExists(ctx, register.Email); err != nil {
		return entity.TokenPair{}, apperror.ErrAuthUserEmailExists.Internal(err)
	}

	if err = s.validateRegister(ctx, register); err != nil {
		return entity.TokenPair{}, err
	}

	if register.Password, err = s.passwordHasher.Hash(ctx, register.Password); err != nil {
		return entity.TokenPair{}, apperror.ErrAuthInternalError.Internal(err)
	}

	err = s.userRepo.Transaction(ctx, func(ctx context.Context) error {
		var session entity.Session
		if session.UserID, err = s.userRepo.Save(ctx, register); err != nil {
			return apperror.ErrAuthDatabaseError.Internal(err)
		}

		if tokens, err = s.sessionManager.Create(ctx, session); err != nil {
			return apperror.ErrAuthInternalError.Internal(err)
		}

		return nil
	})

	if err != nil {
		return entity.TokenPair{}, err
	}

	return tokens, nil
}

func (s *Service) Login(ctx context.Context, creds entity.Credentials) (tokens entity.TokenPair, err error) {
	var user entity.User
	if user, err = s.getUserByEmail(ctx, creds.Email); err != nil {
		return entity.TokenPair{}, err
	}

	if err = s.passwordHasher.Compare(ctx, user.Password, creds.Password); err != nil {
		return entity.TokenPair{}, apperror.ErrAuthUserNotFound.Internal(err)
	}

	if tokens, err = s.sessionManager.Create(ctx, entity.Session{UserID: uint(user.Id)}); err != nil {
		return entity.TokenPair{}, apperror.ErrAuthInternalError.Internal(err)
	}

	return tokens, nil
}

func (s *Service) Refresh(ctx context.Context, refresh entity.RefreshSession) (tokens entity.TokenPair, err error) {
	return s.sessionManager.Refresh(ctx, refresh)
}

func (s *Service) getUserByEmail(ctx context.Context, email string) (user entity.User, err error) {
	if user, err = s.userRepo.GetByEmail(ctx, email); err != nil {
		if errors.Is(err, apperror.ErrAuthUserNotFound) {
			return entity.User{}, nil
		}
		return entity.User{}, apperror.ErrAuthDatabaseError.Internal(err)
	}

	return user, nil
}

func (s *Service) emailExists(ctx context.Context, email string) (err error) {
	_, err = s.userRepo.GetByEmail(ctx, email)
	switch {
	case errors.Is(err, apperror.ErrAuthUserNotFound):
		return nil
	case err != nil:
		return fmt.Errorf("failed to getting from repo user by email %w", err)
	default:
		return fmt.Errorf("user already exists")
	}
}

func (s *Service) validateRegister(ctx context.Context, register entity.Register) (err error) {
	if err = s.validatePassword(ctx, register.Password); err != nil {
		return apperror.ErrAuthPassword.Internal(err)
	}

	if err = s.validateEmail(ctx, register.Email); err != nil {
		return err
	}

	return nil
}

//nolint:cyclop
func (s *Service) validatePassword(_ context.Context, password string) (err error) {
	if !s.config.CheckPassword {
		return nil
	}

	var violations []string

	if len(password) < s.passwordPolicy.MinLength {
		violations = append(violations, fmt.Sprintf("password must be at least %d characters", s.passwordPolicy.MinLength))
	}

	if s.passwordPolicy.RequireLower && !requireLowerRegexp.MatchString(password) {
		violations = append(violations, "must contain at least one lowercase letter")
	}

	if s.passwordPolicy.RequireUpper && !requireUpperRegexp.MatchString(password) {
		violations = append(violations, "must contain at least one uppercase letter")
	}

	if s.passwordPolicy.RequireDigit && !requireDigitRegexp.MatchString(password) {
		violations = append(violations, "must contain at least one digit")
	}

	if len(violations) > 0 {
		return apperror.ErrAuthPassword.SetDetails(violations)
	}

	return nil
}

func (s *Service) validateEmail(ctx context.Context, email string) (err error) {
	if err = s.validateEmailFormat(ctx, email); err != nil {
		return apperror.ErrAuthEmailFormat.Internal(err)
	}

	if s.config.CheckEmailMX {
		if err = s.emailVerifier.VerifyMX(ctx, email); err != nil {
			return apperror.ErrAuthEmailDomain.Internal(err)
		}
	}

	return nil
}

func (s *Service) validateEmailFormat(_ context.Context, email string) (err error) {
	if !s.config.CheckEmailFormat {
		return nil
	}

	if ok := emailFormatRegexp.MatchString(email); !ok {
		return fmt.Errorf("email have invalid format")
	}

	return nil
}

type emailVerifier struct{}

func (emailVerifier) VerifyMX(ctx context.Context, email string) (err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second)
	defer cancelFunc()

	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return fmt.Errorf("malformed email format")
	}

	mx, err := net.DefaultResolver.LookupMX(ctx, parts[1])
	if err != nil {
		return fmt.Errorf("failed to resolve email domain %w", err)
	}

	if len(mx) == 0 {
		return fmt.Errorf("domain doesn't have MX records")
	}

	return nil
}
