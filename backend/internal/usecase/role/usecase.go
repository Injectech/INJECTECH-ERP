package role

import (
	"context"
	"time"

	"github.com/google/uuid"

	domainrole "backend/internal/domain/role"
)

// Service handles role logic.
type Service struct {
	repo domainrole.Repository
	now  func() time.Time
}

func NewService(repo domainrole.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) Create(ctx context.Context, name, description string, permissions []string) (domainrole.Role, error) {
	r := domainrole.Role{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Permissions: permissions,
		CreatedAt:   s.now(),
		UpdatedAt:   s.now(),
	}
	return s.repo.Create(ctx, r)
}

func (s *Service) List(ctx context.Context) ([]domainrole.Role, error) {
	return s.repo.List(ctx)
}
