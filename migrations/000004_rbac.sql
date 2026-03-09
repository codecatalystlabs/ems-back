
-- ============================================
-- File: 000004_rbac.sql
-- ============================================
-- +goose Up
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    module TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(role_id, permission_id)
);

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    scope_type TEXT NOT NULL DEFAULT 'GLOBAL' CHECK (scope_type IN ('GLOBAL', 'DISTRICT', 'FACILITY', 'OWN')),
    scope_id UUID,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    assigned_by UUID REFERENCES users(id),
    UNIQUE(user_id, role_id, scope_type, scope_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_user_roles_scope ON user_roles(scope_type, scope_id);

INSERT INTO roles (code, name, description, is_system) VALUES
('SUPER_ADMIN', 'Super Admin', 'System-wide full access', TRUE),
('NATIONAL_ADMIN', 'National Admin', 'National administrator', TRUE),
('DISTRICT_ADMIN', 'District Admin', 'District-level administrator', TRUE),
('DISPATCH_SUPERVISOR', 'Dispatch Supervisor', 'Dispatch supervisory role', TRUE),
('DISPATCHER', 'Dispatcher', 'Dispatch operator', TRUE),
('CALL_CENTER_AGENT', 'Call Center Agent', 'Call intake operator', TRUE),
('FLEET_MANAGER', 'Fleet Manager', 'Fleet and ambulance manager', TRUE),
('MAINTENANCE_OFFICER', 'Maintenance Officer', 'Maintenance manager', TRUE),
('FACILITY_FOCAL_PERSON', 'Facility Focal Person', 'Facility contact and coordinator', TRUE),
('DRIVER', 'Driver', 'Ambulance driver', TRUE),
('MEDIC', 'Medic', 'Clinical responder', TRUE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO permissions (code, name, module, description) VALUES
('users.create', 'Create users', 'users', 'Can create users'),
('users.read', 'Read users', 'users', 'Can view users'),
('users.update', 'Update users', 'users', 'Can update users'),
('users.deactivate', 'Deactivate users', 'users', 'Can deactivate users'),
('roles.manage', 'Manage roles', 'rbac', 'Can manage roles and permissions'),
('facilities.read', 'Read facilities', 'reference', 'Can view facilities'),
('fleet.read', 'Read fleet', 'fleet', 'Can view ambulances and readiness'),
('fleet.manage', 'Manage fleet', 'fleet', 'Can create and update ambulances'),
('fleet.update_status', 'Update fleet status', 'fleet', 'Can update readiness and status'),
('availability.read', 'Read availability', 'availability', 'Can view user availability'),
('availability.update', 'Update availability', 'availability', 'Can update user availability'),
('incidents.create', 'Create incidents', 'incidents', 'Can create incidents'),
('incidents.read', 'Read incidents', 'incidents', 'Can view incidents'),
('incidents.triage', 'Triage incidents', 'incidents', 'Can triage incidents'),
('incidents.verify', 'Verify incidents', 'incidents', 'Can verify incidents'),
('incidents.escalate', 'Escalate incidents', 'incidents', 'Can escalate incidents'),
('dispatch.read', 'Read dispatch', 'dispatch', 'Can view dispatch assignments'),
('dispatch.assign', 'Assign dispatch', 'dispatch', 'Can assign ambulances and crews'),
('dispatch.update_status', 'Update dispatch status', 'dispatch', 'Can update dispatch lifecycle'),
('trips.read', 'Read trips', 'trips', 'Can view trip data'),
('reports.view', 'View reports', 'reporting', 'Can view reports'),
('audit.read', 'Read audit', 'audit', 'Can view audit logs')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'SUPER_ADMIN'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'users.read', 'users.create', 'users.update', 'users.deactivate',
    'facilities.read',
    'fleet.read', 'fleet.manage', 'fleet.update_status',
    'availability.read', 'availability.update',
    'incidents.read', 'incidents.create', 'incidents.triage', 'incidents.verify', 'incidents.escalate',
    'dispatch.read', 'dispatch.assign', 'dispatch.update_status',
    'trips.read', 'reports.view', 'audit.read'
)
WHERE r.code = 'NATIONAL_ADMIN'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'users.read',
    'facilities.read',
    'fleet.read', 'fleet.manage', 'fleet.update_status',
    'availability.read', 'availability.update',
    'incidents.read', 'incidents.create', 'incidents.triage', 'incidents.verify', 'incidents.escalate',
    'dispatch.read', 'dispatch.assign', 'dispatch.update_status',
    'trips.read', 'reports.view'
)
WHERE r.code = 'DISTRICT_ADMIN'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'facilities.read',
    'fleet.read', 'fleet.update_status',
    'availability.read', 'availability.update',
    'incidents.read', 'incidents.create', 'incidents.triage', 'incidents.verify', 'incidents.escalate',
    'dispatch.read', 'dispatch.assign', 'dispatch.update_status',
    'trips.read', 'reports.view'
)
WHERE r.code = 'DISPATCH_SUPERVISOR'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'facilities.read',
    'fleet.read',
    'availability.read',
    'incidents.read', 'incidents.create', 'incidents.triage', 'incidents.verify',
    'dispatch.read', 'dispatch.assign', 'dispatch.update_status',
    'trips.read'
)
WHERE r.code = 'DISPATCHER'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'incidents.create', 'incidents.read',
    'dispatch.read'
)
WHERE r.code = 'CALL_CENTER_AGENT'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'facilities.read',
    'fleet.read', 'fleet.manage', 'fleet.update_status',
    'availability.read',
    'reports.view'
)
WHERE r.code = 'FLEET_MANAGER'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'fleet.read', 'fleet.update_status'
)
WHERE r.code = 'MAINTENANCE_OFFICER'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'facilities.read',
    'incidents.create', 'incidents.read',
    'dispatch.read',
    'trips.read'
)
WHERE r.code = 'FACILITY_FOCAL_PERSON'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'dispatch.read', 'dispatch.update_status',
    'trips.read'
)
WHERE r.code = 'DRIVER'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'dispatch.read', 'dispatch.update_status',
    'trips.read'
)
WHERE r.code = 'MEDIC'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
