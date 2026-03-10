
-- ============================================
-- File: 000007_incidents.sql
-- ============================================
-- +goose Up
CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_number TEXT NOT NULL UNIQUE,
    source_channel TEXT NOT NULL CHECK (source_channel IN ('SMS', 'USSD', 'CALL', 'MOBILE_APP', 'WEB_PORTAL', 'FACILITY_REFERRAL')),
    caller_name TEXT,
    caller_phone TEXT,
    patient_name TEXT,
    patient_phone TEXT,
    patient_age_group TEXT,
    patient_sex TEXT CHECK (patient_sex IN ('MALE', 'FEMALE', 'OTHER', 'UNKNOWN')),
    incident_type_id UUID NOT NULL REFERENCES ref_incident_types(id),
    severity_level_id UUID REFERENCES ref_severity_levels(id),
    priority_level_id UUID REFERENCES ref_priority_levels(id),
    summary TEXT,
    description TEXT,
    district_id UUID REFERENCES ref_districts(id),
    facility_id UUID REFERENCES ref_facilities(id),
    village TEXT,
    parish TEXT,
    subcounty TEXT,
    landmark TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    location GEOGRAPHY(POINT, 4326),
    verification_status TEXT NOT NULL DEFAULT 'PENDING' CHECK (verification_status IN ('PENDING', 'VERIFIED', 'REJECTED')),
    status TEXT NOT NULL DEFAULT 'NEW' CHECK (status IN ('NEW', 'PENDING_VERIFICATION', 'VERIFIED', 'AWAITING_ASSIGNMENT', 'ASSIGNED', 'ENROUTE', 'AT_SCENE', 'TRANSPORTING', 'COMPLETED', 'CANCELLED', 'ESCALATED', 'REJECTED')),
    reported_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by_user_id UUID REFERENCES users(id),
    triaged_by_user_id UUID REFERENCES users(id),
    triaged_at TIMESTAMPTZ,
    assigned_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_incidents_incident_type_id ON incidents(incident_type_id);
CREATE INDEX idx_incidents_priority_level_id ON incidents(priority_level_id);
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_verification_status ON incidents(verification_status);
CREATE INDEX idx_incidents_district_id ON incidents(district_id);
CREATE INDEX idx_incidents_facility_id ON incidents(facility_id);
CREATE INDEX idx_incidents_reported_at ON incidents(reported_at DESC);
CREATE INDEX idx_incidents_location ON incidents USING GIST(location);

CREATE TABLE incident_updates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    update_type TEXT NOT NULL CHECK (update_type IN ('COMMENT', 'STATUS_CHANGE', 'TRIAGE', 'VERIFICATION', 'ESCALATION', 'LOCATION_UPDATE', 'CANCELLATION')),
    old_value TEXT,
    new_value TEXT,
    notes TEXT,
    actor_user_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_incident_updates_incident_id ON incident_updates(incident_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS incident_updates;
DROP TABLE IF EXISTS incidents;
