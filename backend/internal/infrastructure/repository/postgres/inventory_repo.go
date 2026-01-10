package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domaininventory "backend/internal/domain/inventory"
	"backend/internal/infrastructure/repository/postgres/sqlc"
)

// InventoryRepository implements inventory.Repository via PostgreSQL.
type InventoryRepository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{db: db, q: sqlc.New(db)}
}

func (r *InventoryRepository) Create(ctx context.Context, inv domaininventory.Inventory) (domaininventory.Inventory, error) {
	id, err := uuid.Parse(inv.ID)
	if err != nil {
		return domaininventory.Inventory{}, err
	}
	pid, err := uuid.Parse(inv.ProductID)
	if err != nil {
		return domaininventory.Inventory{}, err
	}
	created, err := r.q.CreateInventory(ctx, sqlc.CreateInventoryParams{
		ID:        id,
		ProductID: pid,
		Quantity:  inv.Quantity,
		Location:  inv.Location,
		CreatedAt: inv.CreatedAt,
		UpdatedAt: inv.UpdatedAt,
	})
	if err != nil {
		return domaininventory.Inventory{}, err
	}
	inv.CreatedAt = created.CreatedAt
	inv.UpdatedAt = created.UpdatedAt
	return inv, nil
}

func (r *InventoryRepository) GetByID(ctx context.Context, id string) (domaininventory.Inventory, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domaininventory.Inventory{}, err
	}
	inv, err := r.q.GetInventoryByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domaininventory.Inventory{}, err
		}
		return domaininventory.Inventory{}, err
	}
	return domaininventory.Inventory{
		ID:        inv.ID.String(),
		ProductID: inv.ProductID.String(),
		Quantity:  inv.Quantity,
		Location:  inv.Location,
		CreatedAt: inv.CreatedAt,
		UpdatedAt: inv.UpdatedAt,
		DeletedAt: inv.DeletedAt,
	}, nil
}

func (r *InventoryRepository) ListByProduct(ctx context.Context, productID string) ([]domaininventory.Inventory, error) {
	pid, err := uuid.Parse(productID)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListInventoryByProduct(ctx, pid)
	if err != nil {
		return nil, err
	}

	var res []domaininventory.Inventory
	for _, inv := range rows {
		res = append(res, domaininventory.Inventory{
			ID:        inv.ID.String(),
			ProductID: inv.ProductID.String(),
			Quantity:  inv.Quantity,
			Location:  inv.Location,
			CreatedAt: inv.CreatedAt,
			UpdatedAt: inv.UpdatedAt,
			DeletedAt: inv.DeletedAt,
		})
	}
	return res, nil
}

func (r *InventoryRepository) List(ctx context.Context) ([]domaininventory.Inventory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, product_id, quantity, location, created_at, updated_at, deleted_at
		FROM inventories
		WHERE deleted_at IS NULL
		ORDER BY location ASC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domaininventory.Inventory
	for rows.Next() {
		var inv domaininventory.Inventory
		var id uuid.UUID
		var pid uuid.UUID
		if err := rows.Scan(
			&id,
			&pid,
			&inv.Quantity,
			&inv.Location,
			&inv.CreatedAt,
			&inv.UpdatedAt,
			&inv.DeletedAt,
		); err != nil {
			return nil, err
		}
		inv.ID = id.String()
		inv.ProductID = pid.String()
		res = append(res, inv)
	}
	return res, rows.Err()
}

func (r *InventoryRepository) Adjust(ctx context.Context, id string, delta int64) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.AdjustInventory(ctx, sqlc.AdjustInventoryParams{
		ID:        uid,
		Delta:     delta,
		UpdatedAt: time.Now(),
	})
}

func (r *InventoryRepository) UpdateLocation(ctx context.Context, id, location string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `
		UPDATE inventories
		SET location = $2, updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL
	`, uid, location, time.Now())
	return err
}

func (r *InventoryRepository) SoftDelete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.SoftDeleteInventory(ctx, uid)
}
