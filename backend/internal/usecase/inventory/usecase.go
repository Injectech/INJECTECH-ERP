package inventory

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	domaininventory "backend/internal/domain/inventory"
	domainlocation "backend/internal/domain/location"
)

// Service manages inventory adjustments.
type Service struct {
	repo      domaininventory.Repository
	locations domainlocation.Repository
}

func NewService(repo domaininventory.Repository, locations domainlocation.Repository) *Service {
	return &Service{repo: repo, locations: locations}
}

func (s *Service) Create(ctx context.Context, inv domaininventory.Inventory) (domaininventory.Inventory, error) {
	if inv.ID == "" {
		inv.ID = uuid.NewString()
	}
	if inv.CreatedAt.IsZero() {
		inv.CreatedAt = time.Now()
	}
	if inv.UpdatedAt.IsZero() {
		inv.UpdatedAt = inv.CreatedAt
	}
	if inv.Location == "" && s.locations != nil {
		if loc, err := s.locations.GetDefault(ctx); err == nil {
			inv.Location = loc.Name
		} else if errors.Is(err, domainlocation.ErrNotFound) {
			inv.Location = "Default"
		} else {
			return domaininventory.Inventory{}, err
		}
	}
	return s.repo.Create(ctx, inv)
}

func (s *Service) Adjust(ctx context.Context, id string, delta int64) error {
	return s.repo.Adjust(ctx, id, delta)
}

func (s *Service) UpdateLocation(ctx context.Context, id, location string) error {
	loc := location
	if loc == "" && s.locations != nil {
		if def, err := s.locations.GetDefault(ctx); err == nil {
			loc = def.Name
		} else if errors.Is(err, domainlocation.ErrNotFound) {
			loc = "Default"
		} else {
			return err
		}
	}
	if loc == "" {
		loc = "Default"
	}
	return s.repo.UpdateLocation(ctx, id, loc)
}

func (s *Service) List(ctx context.Context) ([]domaininventory.Inventory, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListByProduct(ctx context.Context, productID string) ([]domaininventory.Inventory, error) {
	return s.repo.ListByProduct(ctx, productID)
}
