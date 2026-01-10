-- name: CreateInventory :one
INSERT INTO inventories (id, product_id, quantity, location, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetInventoryByID :one
SELECT id, product_id, quantity, location, created_at, updated_at, deleted_at
FROM inventories
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListInventoryByProduct :many
SELECT id, product_id, quantity, location, created_at, updated_at, deleted_at
FROM inventories
WHERE product_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: AdjustInventory :exec
UPDATE inventories
SET quantity = quantity + $2, updated_at = $3
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteInventory :exec
UPDATE inventories SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL;
