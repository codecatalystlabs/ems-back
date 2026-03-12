-- ============================================
-- File: 000026_rbac_permissions_extra_sql.sql
-- ============================================
-- +goose Up
INSERT INTO permissions (code, name, module, description) VALUES
('facilities.manage', 'Manage facilities', 'reference', 'Can create/update/delete facilities'),
('trips.manage', 'Manage trips', 'trips', 'Can create/update/delete trips and events'),
('notifications.read', 'Read notifications', 'notifications', 'Can view notifications'),
('notifications.manage', 'Manage notifications', 'notifications', 'Can create/mark notifications'),
('fuel.read', 'Read fuel logs', 'fuel', 'Can view fuel logs'),
('fuel.manage', 'Manage fuel logs', 'fuel', 'Can create/update/delete fuel logs')
ON CONFLICT (code) DO NOTHING;

-- SUPER_ADMIN already gets all permissions via the CROSS JOIN in the original migration,
-- but we still grant explicitly for safety in case environments skipped that step.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
  'facilities.manage',
  'trips.manage',
  'notifications.read',
  'notifications.manage',
  'fuel.read',
  'fuel.manage'
)
WHERE r.code = 'SUPER_ADMIN'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
  'facilities.manage',
  'trips.manage',
  'notifications.read',
  'notifications.manage',
  'fuel.read',
  'fuel.manage'
)
WHERE r.code IN ('NATIONAL_ADMIN', 'DISTRICT_ADMIN', 'FLEET_MANAGER', 'DISPATCH_SUPERVISOR')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- +goose Down
DELETE FROM role_permissions rp
USING roles r, permissions p
WHERE rp.role_id = r.id
  AND rp.permission_id = p.id
  AND p.code IN (
    'facilities.manage',
    'trips.manage',
    'notifications.read',
    'notifications.manage',
    'fuel.read',
    'fuel.manage'
  );

DELETE FROM permissions
WHERE code IN (
  'facilities.manage',
  'trips.manage',
  'notifications.read',
  'notifications.manage',
  'fuel.read',
  'fuel.manage'
);

