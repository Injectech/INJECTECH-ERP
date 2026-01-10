DELETE FROM role_permissions rp
USING roles r, permissions p
WHERE rp.role_id = r.id
  AND rp.permission_id = p.id
  AND r.name = 'admin'
  AND p.code IN ('location.read', 'location.write');

DELETE FROM permissions WHERE code IN ('location.read', 'location.write');
