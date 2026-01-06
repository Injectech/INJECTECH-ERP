package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	domaininventory "backend/internal/domain/inventory"
)

// InventoryRepository implements inventory.Repository via PostgreSQL.
type InventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(ctx context.Context, inv domaininventory.Inventory) (domaininventory.Inventory, error) {
	return inv, nil
}

func (r *InventoryRepository) GetByID(ctx context.Context, id string) (domaininventory.Inventory, error) {
	return domaininventory.Inventory{}, nil
}

func (r *InventoryRepository) ListByProduct(ctx context.Context, productID string) ([]domaininventory.Inventory, error) {
	return []domaininventory.Inventory{}, nil
}

func (r *InventoryRepository) Adjust(ctx context.Context, id string, delta int64) error {
	return nil
}

func (r *InventoryRepository) SoftDelete(ctx context.Context, id string) error {
	return nil
}
