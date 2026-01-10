package auth

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	domainauth "backend/internal/domain/auth"
	domainuser "backend/internal/domain/user"
	"backend/internal/pkg/security"
)

// Service handles authentication workflows.
type Service struct {
	users         domainuser.Repository
	accessTTL     time.Duration
	refreshTTL    time.Duration
	accessSecret  string
	refreshSecret string
	now           func() time.Time
	refreshStore  map[string]string
	storeMu       sync.Mutex
}

// Session represents an authenticated session with tokens and user info.
type Session struct {
	Tokens domainauth.TokenPair
	User   domainuser.User
}

var (
	ErrInvalidCredential = errors.New("invalid credential")
	ErrInvalidRefresh    = errors.New("invalid refresh token")
)

func NewService(users domainuser.Repository, accessTTL, refreshTTL time.Duration, accessSecret, refreshSecret string) *Service {
	return &Service{
		users:         users,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		now:           time.Now,
		refreshStore:  make(map[string]string),
	}
}

// Register creates a new user and returns an authenticated session.
func (s *Service) Register(ctx context.Context, email, password, name string) (Session, error) {
	hashed, err := security.HashPassword(password)
	if err != nil {
		return Session{}, err
	}

	user := domainuser.User{
		ID:             uuid.NewString(),
		Email:          email,
		HashedPassword: hashed,
		Name:           name,
		Roles:          []string{"user"},
		CreatedAt:      s.now(),
		UpdatedAt:      s.now(),
	}

	created, err := s.users.Create(ctx, user)
	if err != nil {
		return Session{}, err
	}

	tokens, err := s.issueTokens(created.ID, created.Roles, created.Permissions)
	if err != nil {
		return Session{}, err
	}

	return Session{Tokens: tokens, User: created}, nil
}

// Login verifies credential and returns an authenticated session.
func (s *Service) Login(ctx context.Context, email, password string) (Session, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return Session{}, err
	}
	if err := security.ComparePassword(u.HashedPassword, password); err != nil {
		return Session{}, ErrInvalidCredential
	}

	tokens, err := s.issueTokens(u.ID, u.Roles, u.Permissions)
	if err != nil {
		return Session{}, err
	}

	return Session{Tokens: tokens, User: u}, nil
}

func (s *Service) issueTokens(subject string, roles, permissions []string) (domainauth.TokenPair, error) {
	now := s.now()
	accessToken, accessExp, err := security.SignToken(s.accessSecret, subject, roles, permissions, s.accessTTL, now)
	if err != nil {
		return domainauth.TokenPair{}, err
	}

	refreshToken, refreshExp, err := security.SignToken(s.refreshSecret, subject, roles, permissions, s.refreshTTL, now)
	if err != nil {
		return domainauth.TokenPair{}, err
	}

	tokens := domainauth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessExp:    accessExp,
		RefreshExp:   refreshExp,
	}
	s.storeRefresh(subject, refreshToken)
	return tokens, nil
}

// ValidateRefresh validates refresh token and returns subject.
func (s *Service) ValidateRefresh(refresh string) (*security.Claims, error) {
	claims, err := security.ParseToken(s.refreshSecret, refresh)
	if err != nil {
		return nil, err
	}
	if !s.isStoredRefresh(claims.Subject, refresh) {
		return nil, ErrInvalidRefresh
	}
	return claims, nil
}

// ParseAccess validates access token.
func (s *Service) ParseAccess(access string) (*security.Claims, error) {
	claims, err := security.ParseToken(s.accessSecret, access)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// Refresh issues a new token pair based on refresh token claims.
func (s *Service) Refresh(refresh string) (domainauth.TokenPair, error) {
	claims, err := s.ValidateRefresh(refresh)
	if err != nil {
		return domainauth.TokenPair{}, err
	}
	return s.issueTokens(claims.Subject, claims.Roles, claims.Permissions)
}

// IsTokenExpired checks if the error indicates expiration.
func IsTokenExpired(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidClaims)
}

func (s *Service) storeRefresh(subject, token string) {
	s.storeMu.Lock()
	defer s.storeMu.Unlock()
	s.refreshStore[subject] = token
}

func (s *Service) isStoredRefresh(subject, token string) bool {
	s.storeMu.Lock()
	defer s.storeMu.Unlock()
	val, ok := s.refreshStore[subject]
	return ok && val == token
}
