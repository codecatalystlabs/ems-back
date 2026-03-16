-- ============================================
-- File: 000029_default_admin_user.sql
-- ============================================
-- +goose Up
-- Default admin user: username=admin, password=admin123
-- Change password after first login.
INSERT INTO users (
    id, username, first_name, last_name, email, phone,
    password_hash, status, is_active, created_at, updated_at
)
SELECT
    gen_random_uuid(),
    'admin',
    'System',
    'Administrator',
    'admin@dispatch.local',
    '+256780000000',
    crypt('admin123', gen_salt('bf', 10)),
    'ACTIVE',
    true,
    now(),
    now()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin' AND deleted_at IS NULL);

-- Assign SUPER_ADMIN role to admin user
INSERT INTO user_roles (user_id, role_id, scope_type, active, assigned_at)
SELECT u.id, r.id, 'GLOBAL', true, now()
FROM users u
CROSS JOIN roles r
WHERE u.username = 'admin' AND u.deleted_at IS NULL
  AND r.code = 'SUPER_ADMIN'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur
    WHERE ur.user_id = u.id AND ur.role_id = r.id AND ur.scope_type = 'GLOBAL'
  );

-- +goose Down
-- Remove default admin (only if still using default email)
DELETE FROM user_roles
WHERE user_id IN (SELECT id FROM users WHERE username = 'admin' AND email = 'admin@dispatch.local');

DELETE FROM users
WHERE username = 'admin' AND email = 'admin@dispatch.local';
