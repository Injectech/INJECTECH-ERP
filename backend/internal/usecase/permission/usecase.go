package permission

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainpermission "backend/internal/domain/permission"
)

// Service manages permissions.
type Service struct {
    repo domainpermission.Repository
    now  func() time.Time
}

func NewService(repo domainpermission.Repository) *Service {
    return &Service{repo: repo, now: time.Now}
}

func (s *Service) Create(ctx context.Context, code, description string) (domainpermission.Permission, error) {
    p := domainpermission.Permission{
        ID:          uuid.NewString(),
        Code:        code,
        Description: description,
        CreatedAt:   s.now(),
        UpdatedAt:   s.now(),
    }
    return s.repo.Create(ctx, p)
}

func (s *Service) List(ctx context.Context) ([]domainpermission.Permission, error) {
    return s.repo.List(ctx)
}
