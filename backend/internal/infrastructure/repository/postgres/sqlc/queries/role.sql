-- name: CreateRole :one
INSERT INTO roles (id, name, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRoleByID :one
SELECT r.id, r.name, r.description, r.created_at, r.updated_at, r.deleted_at
FROM roles r
WHERE r.id = $1 AND r.deleted_at IS NULL;

-- name: GetRoleByName :one
SELECT r.id, r.name, r.description, r.created_at, r.updated_at, r.deleted_at
FROM roles r
WHERE r.name = $1 AND r.deleted_at IS NULL;

-- name: ListRoles :many
SELECT r.id, r.name, r.description, r.created_at, r.updated_at, r.deleted_at
FROM roles r
WHERE r.deleted_at IS NULL
ORDER BY r.created_at DESC;

-- name: SoftDeleteRole :exec
UPDATE roles SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateRole :exec
UPDATE roles
SET name = $2, description = $3, updated_at = $4
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpsertRolePermission :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, $2)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- name: ListRolePermissions :many
SELECT p.code
FROM role_permissions rp
JOIN permissions p ON p.id = rp.permission_id
WHERE rp.role_id = $1;
