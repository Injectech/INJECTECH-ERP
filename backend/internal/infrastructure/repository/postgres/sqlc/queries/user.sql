-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password, name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT u.id, u.email, u.hashed_password, u.name, u.created_at, u.updated_at, u.deleted_at
FROM users u
WHERE u.id = $1 AND u.deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT u.id, u.email, u.hashed_password, u.name, u.created_at, u.updated_at, u.deleted_at
FROM users u
WHERE u.email = $1 AND u.deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUser :exec
UPDATE users
SET email = $2, hashed_password = $3, name = $4, updated_at = $5
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUserRoles :many
SELECT r.name
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: ListUserPermissions :many
SELECT DISTINCT p.code
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1;
