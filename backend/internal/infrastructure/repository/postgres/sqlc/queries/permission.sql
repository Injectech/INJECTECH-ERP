-- name: CreatePermission :one
INSERT INTO permissions (id, code, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPermissionByCode :one
SELECT id, code, description, created_at, updated_at, deleted_at
FROM permissions
WHERE code = $1 AND deleted_at IS NULL;

-- name: ListPermissions :many
SELECT id, code, description, created_at, updated_at, deleted_at
FROM permissions
WHERE deleted_at IS NULL
ORDER BY created_at DESC;
