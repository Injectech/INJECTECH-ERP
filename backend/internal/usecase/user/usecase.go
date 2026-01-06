package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	domainuser "backend/internal/domain/user"
	"backend/internal/pkg/security"
)

// Service contains user business logic.
type Service struct {
	repo domainuser.Repository
	now  func() time.Time
}

var (
	ErrNotFound = errors.New("user not found")
)

func NewService(repo domainuser.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) Create(ctx context.Context, email, password, name string, roles []string) (domainuser.User, error) {
	hashed, err := security.HashPassword(password)
	if err != nil {
		return domainuser.User{}, err
	}
	u := domainuser.User{
		ID:             uuid.NewString(),
		Email:          email,
		HashedPassword: hashed,
		Name:           name,
		Roles:          roles,
		CreatedAt:      s.now(),
		UpdatedAt:      s.now(),
	}
	return s.repo.Create(ctx, u)
}

func (s *Service) Get(ctx context.Context, id string) (domainuser.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domainuser.User{}, err
	}
	return u, nil
}

func (s *Service) Update(ctx context.Context, u domainuser.User) error {
	u.UpdatedAt = s.now()
	return s.repo.Update(ctx, u)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
