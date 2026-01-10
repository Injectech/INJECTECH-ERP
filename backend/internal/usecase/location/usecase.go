package location

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	domainlocation "backend/internal/domain/location"
)

// Service contains location business logic.
type Service struct {
	repo domainlocation.Repository
	now  func() time.Time
}

func NewService(repo domainlocation.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) Create(ctx context.Context, name, description string) (domainlocation.Location, error) {
	isDefault := false
	if _, err := s.repo.GetDefault(ctx); err != nil {
		if errors.Is(err, domainlocation.ErrNotFound) {
			isDefault = true
		} else {
			return domainlocation.Location{}, err
		}
	}

	loc := domainlocation.Location{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		IsDefault:   isDefault,
		CreatedAt:   s.now(),
		UpdatedAt:   s.now(),
	}
	return s.repo.Create(ctx, loc)
}

func (s *Service) List(ctx context.Context) ([]domainlocation.Location, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetDefault(ctx context.Context) (domainlocation.Location, error) {
	return s.repo.GetDefault(ctx)
}
