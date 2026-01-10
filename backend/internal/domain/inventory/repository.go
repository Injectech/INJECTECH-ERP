package inventory

import "context"

// Repository defines persistence operations for inventory.
type Repository interface {
	Create(ctx context.Context, inv Inventory) (Inventory, error)
	GetByID(ctx context.Context, id string) (Inventory, error)
	List(ctx context.Context) ([]Inventory, error)
	ListByProduct(ctx context.Context, productID string) ([]Inventory, error)
	Adjust(ctx context.Context, id string, delta int64) error
	UpdateLocation(ctx context.Context, id, location string) error
	SoftDelete(ctx context.Context, id string) error
}
