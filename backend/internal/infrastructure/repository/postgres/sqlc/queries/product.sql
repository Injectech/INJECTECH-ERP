-- name: CreateProduct :one
INSERT INTO products (id, sku, name, description, price, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetProductByID :one
SELECT id, sku, name, description, price, created_at, updated_at, deleted_at
FROM products
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListProducts :many
SELECT id, sku, name, description, price, created_at, updated_at, deleted_at
FROM products
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateProduct :exec
UPDATE products
SET sku = $2, name = $3, description = $4, price = $5, updated_at = $6
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteProduct :exec
UPDATE products SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL;
