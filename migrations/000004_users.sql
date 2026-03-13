-- ============================================
-- File: 000004_users.sql
-- ============================================
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_no TEXT UNIQUE,
    username CITEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    other_name TEXT,
    gender TEXT,
    phone TEXT UNIQUE,
    email CITEXT UNIQUE,
    password_hash TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED', 'LOCKED')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_locked BOOLEAN NOT NULL DEFAULT FALSE,
    failed_login_attempts INT NOT NULL DEFAULT 0 CHECK (failed_login_attempts >= 0),
    last_login_at TIMESTAMPTZ,
    password_changed_at TIMESTAMPTZ,
    preferred_language TEXT NOT NULL DEFAULT 'en',
    timezone TEXT NOT NULL DEFAULT 'Africa/Kampala',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_is_active ON users(is_active);

CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    cadre TEXT,
    license_number TEXT,
    specialization TEXT,
    date_of_birth DATE,
    national_id TEXT,
    avatar_url TEXT,
    emergency_contact_name TEXT,
    emergency_contact_phone TEXT,
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    district_id UUID REFERENCES ref_districts(id),
    subcounty_id UUID REFERENCES ref_subcounties(id),
    facility_id UUID REFERENCES ref_facilities(id),
    assignment_level TEXT NOT NULL CHECK (
        assignment_level IN ('NATIONAL', 'DISTRICT', 'SUBCOUNTY', 'FACILITY')
    ),
    team_name TEXT,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_user_assignments_dates CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE INDEX idx_user_assignments_user_id ON user_assignments(user_id);
CREATE INDEX idx_user_assignments_district_id ON user_assignments(district_id);
CREATE INDEX idx_user_assignments_facility_id ON user_assignments(facility_id);
CREATE UNIQUE INDEX uq_user_primary_assignment ON user_assignments(user_id) WHERE is_primary = TRUE AND active = TRUE;

CREATE TABLE user_capabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    capability_id UUID NOT NULL REFERENCES ref_capabilities(id),
    level_no INT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, capability_id)
);

CREATE INDEX idx_user_capabilities_user_id ON user_capabilities(user_id);
CREATE INDEX idx_user_capabilities_capability_id ON user_capabilities(capability_id);

-- +goose Down
DROP TABLE IF EXISTS user_capabilities;
DROP TABLE IF EXISTS user_assignments;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS users;

