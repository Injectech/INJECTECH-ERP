INSERT INTO permissions (id, code, description, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'location.read', 'Read locations', now(), now()),
    (uuid_generate_v4(), 'location.write', 'Create/Update locations', now(), now())
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN ('location.read', 'location.write')
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;
