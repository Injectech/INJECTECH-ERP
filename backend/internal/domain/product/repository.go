package product

import "context"

// Repository defines persistence operations for products.
type Repository interface {
	Create(ctx context.Context, p Product) (Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, p Product) error
	SoftDelete(ctx context.Context, id string) error
}
