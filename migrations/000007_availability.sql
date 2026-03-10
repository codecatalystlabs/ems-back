-- ============================================
-- File: 000006_availability.sql
-- ============================================
-- +goose Up
CREATE TABLE user_shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shift_date DATE NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    shift_type TEXT,
    district_id UUID REFERENCES ref_districts(id),
    facility_id UUID REFERENCES ref_facilities(id),
    status TEXT NOT NULL DEFAULT 'SCHEDULED' CHECK (status IN ('SCHEDULED', 'ACTIVE', 'COMPLETED', 'CANCELLED', 'MISSED')),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_user_shifts_range CHECK (ends_at > starts_at)
);

CREATE INDEX idx_user_shifts_user_id ON user_shifts(user_id);
CREATE INDEX idx_user_shifts_date ON user_shifts(shift_date);
CREATE INDEX idx_user_shifts_status ON user_shifts(status);

CREATE TABLE user_availability (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    availability_status TEXT NOT NULL CHECK (availability_status IN ('AVAILABLE', 'BUSY', 'OFF_DUTY', 'ON_LEAVE', 'STANDBY', 'OFFLINE', 'SUSPENDED')),
    dispatchable BOOLEAN NOT NULL DEFAULT FALSE,
    current_incident_id UUID,
    current_dispatch_assignment_id UUID,
    current_ambulance_id UUID REFERENCES ambulances(id),
    last_seen_at TIMESTAMPTZ,
    source TEXT NOT NULL DEFAULT 'SYSTEM' CHECK (source IN ('SYSTEM', 'MANUAL', 'APP', 'SYNC', 'SHIFT_ENGINE')),
    notes TEXT,
    updated_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_availability_status ON user_availability(availability_status);
CREATE INDEX idx_user_availability_dispatchable ON user_availability(dispatchable);

CREATE TABLE user_presence_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel TEXT NOT NULL CHECK (channel IN ('WEB', 'MOBILE', 'SMS', 'USSD', 'CALL_CENTER', 'SYSTEM')),
    seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ip_address INET,
    user_agent TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION
);

CREATE INDEX idx_user_presence_logs_user_id ON user_presence_logs(user_id, seen_at DESC);

-- +goose Down
DROP TABLE IF EXISTS user_presence_logs;
DROP TABLE IF EXISTS user_availability;
DROP TABLE IF EXISTS user_shifts;
