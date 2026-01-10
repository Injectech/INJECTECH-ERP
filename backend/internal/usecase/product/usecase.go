package product

import (
	"context"
	"time"

	"github.com/google/uuid"

	domainproduct "backend/internal/domain/product"
)

// Service contains product business logic.
type Service struct {
	repo domainproduct.Repository
	now  func() time.Time
}

func NewService(repo domainproduct.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) Create(ctx context.Context, sku, name, description string, price float64) (domainproduct.Product, error) {
	p := domainproduct.Product{
		ID:          uuid.NewString(),
		SKU:         sku,
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   s.now(),
		UpdatedAt:   s.now(),
	}
	return s.repo.Create(ctx, p)
}

func (s *Service) List(ctx context.Context) ([]domainproduct.Product, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id, sku, name, description string, price float64) (domainproduct.Product, error) {
	p := domainproduct.Product{
		ID:          id,
		SKU:         sku,
		Name:        name,
		Description: description,
		Price:       price,
		UpdatedAt:   s.now(),
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return domainproduct.Product{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
