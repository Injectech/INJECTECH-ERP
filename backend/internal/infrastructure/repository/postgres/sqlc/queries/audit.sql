-- name: CreateAuditLog :exec
INSERT INTO audit_logs (id, actor_id, action, resource, metadata, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListAuditLogs :many
SELECT id, actor_id, action, resource, metadata, created_at
FROM audit_logs
WHERE ($1::uuid IS NULL OR actor_id = $1)
ORDER BY created_at DESC;
