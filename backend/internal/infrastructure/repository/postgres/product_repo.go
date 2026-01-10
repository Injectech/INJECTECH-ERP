package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainproduct "backend/internal/domain/product"
	"backend/internal/infrastructure/repository/postgres/sqlc"
)

// ProductRepository implements product.Repository via PostgreSQL.
type ProductRepository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db, q: sqlc.New(db)}
}

func (r *ProductRepository) Create(ctx context.Context, p domainproduct.Product) (domainproduct.Product, error) {
	pid, err := uuid.Parse(p.ID)
	if err != nil {
		return domainproduct.Product{}, err
	}
	created, err := r.q.CreateProduct(ctx, sqlc.CreateProductParams{
		ID:          pid,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	})
	if err != nil {
		return domainproduct.Product{}, err
	}
	p.CreatedAt = created.CreatedAt
	p.UpdatedAt = created.UpdatedAt
	return p, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (domainproduct.Product, error) {
	pid, err := uuid.Parse(id)
	if err != nil {
		return domainproduct.Product{}, err
	}
	p, err := r.q.GetProductByID(ctx, pid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainproduct.Product{}, err
		}
		return domainproduct.Product{}, err
	}
	return domainproduct.Product{
		ID:          p.ID.String(),
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt,
	}, nil
}

func (r *ProductRepository) List(ctx context.Context) ([]domainproduct.Product, error) {
	dbProducts, err := r.q.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	var products []domainproduct.Product
	for _, p := range dbProducts {
		products = append(products, domainproduct.Product{
			ID:          p.ID.String(),
			SKU:         p.SKU,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			DeletedAt:   p.DeletedAt,
		})
	}
	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, p domainproduct.Product) error {
	pid, err := uuid.Parse(p.ID)
	if err != nil {
		return err
	}
	return r.q.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:          pid,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		UpdatedAt:   time.Now(),
	})
}

func (r *ProductRepository) SoftDelete(ctx context.Context, id string) error {
	pid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.SoftDeleteProduct(ctx, pid)
}
