package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	domainproduct "backend/internal/domain/product"
)

// ProductRepository implements product.Repository via PostgreSQL.
type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p domainproduct.Product) (domainproduct.Product, error) {
	return p, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (domainproduct.Product, error) {
	return domainproduct.Product{}, nil
}

func (r *ProductRepository) List(ctx context.Context) ([]domainproduct.Product, error) {
	return []domainproduct.Product{}, nil
}

func (r *ProductRepository) Update(ctx context.Context, p domainproduct.Product) error {
	return nil
}

func (r *ProductRepository) SoftDelete(ctx context.Context, id string) error {
	return nil
}
