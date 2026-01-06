package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"

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
}

var (
	ErrInvalidCredential = errors.New("invalid credential")
)

func NewService(users domainuser.Repository, accessTTL, refreshTTL time.Duration, accessSecret, refreshSecret string) *Service {
	return &Service{
		users:         users,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		now:           time.Now,
	}
}

// Register creates a new user and returns a token pair.
func (s *Service) Register(ctx context.Context, email, password, name string) (domainauth.TokenPair, error) {
	hashed, err := security.HashPassword(password)
	if err != nil {
		return domainauth.TokenPair{}, err
	}

	user := domainuser.User{
		ID:        uuid.NewString(),
		Email:     email,
		HashedPassword: hashed,
		Name:      name,
		Roles:     []string{"user"},
		CreatedAt: s.now(),
		UpdatedAt: s.now(),
	}

	created, err := s.users.Create(ctx, user)
	if err != nil {
		return domainauth.TokenPair{}, err
	}

	return s.issueTokens(created.ID, created.Roles)
}

// Login verifies credential and returns a token pair.
func (s *Service) Login(ctx context.Context, email, password string) (domainauth.TokenPair, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return domainauth.TokenPair{}, err
	}
	if err := security.ComparePassword(u.HashedPassword, password); err != nil {
		return domainauth.TokenPair{}, ErrInvalidCredential
	}

	return s.issueTokens(u.ID, u.Roles)
}

func (s *Service) issueTokens(subject string, roles []string) (domainauth.TokenPair, error) {
	accessToken, accessExp, err := security.SignToken(s.accessSecret, subject, roles, s.accessTTL)
	if err != nil {
		return domainauth.TokenPair{}, err
	}

	refreshToken, refreshExp, err := security.SignToken(s.refreshSecret, subject, roles, s.refreshTTL)
	if err != nil {
		return domainauth.TokenPair{}, err
	}

	return domainauth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessExp:    accessExp,
		RefreshExp:   refreshExp,
	}, nil
}

// ValidateRefresh validates refresh token and returns subject.
func (s *Service) ValidateRefresh(refresh string) (*security.Claims, error) {
	claims, err := security.ParseToken(s.refreshSecret, refresh)
	if err != nil {
		return nil, err
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
	return s.issueTokens(claims.Subject, claims.Roles)
}

// IsTokenExpired checks if the error indicates expiration.
func IsTokenExpired(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidClaims)
}
