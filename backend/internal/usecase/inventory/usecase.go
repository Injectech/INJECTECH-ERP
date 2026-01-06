package inventory

import (
	"context"

	domaininventory "backend/internal/domain/inventory"
)

// Service manages inventory adjustments.
type Service struct {
	repo domaininventory.Repository
}

func NewService(repo domaininventory.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, inv domaininventory.Inventory) (domaininventory.Inventory, error) {
	return s.repo.Create(ctx, inv)
}

func (s *Service) Adjust(ctx context.Context, id string, delta int64) error {
	return s.repo.Adjust(ctx, id, delta)
}

func (s *Service) ListByProduct(ctx context.Context, productID string) ([]domaininventory.Inventory, error) {
	return s.repo.ListByProduct(ctx, productID)
}
