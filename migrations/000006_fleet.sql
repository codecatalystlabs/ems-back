
-- ============================================
-- File: 000006_fleet.sql
-- ============================================
-- +goose Up
CREATE TABLE ambulances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT UNIQUE,
    plate_number TEXT NOT NULL UNIQUE,
    vin TEXT UNIQUE,
    make TEXT,
    model TEXT,
    year_of_manufacture INT,
    category_id UUID NOT NULL REFERENCES ref_ambulance_categories(id),
    ownership_type TEXT,
    station_facility_id UUID REFERENCES ref_facilities(id),
    district_id UUID REFERENCES ref_districts(id),
    status TEXT NOT NULL DEFAULT 'AVAILABLE' CHECK (status IN ('AVAILABLE', 'RESERVED', 'ASSIGNED', 'ENROUTE', 'AT_SCENE', 'TRANSPORTING', 'RETURNING', 'MAINTENANCE', 'BREAKDOWN', 'OFFLINE', 'RETIRED')),
    dispatch_readiness TEXT NOT NULL DEFAULT 'DISPATCHABLE' CHECK (dispatch_readiness IN ('DISPATCHABLE', 'RESTRICTED', 'NOT_DISPATCHABLE')),
    gps_lat DOUBLE PRECISION,
    gps_lon DOUBLE PRECISION,
    location GEOGRAPHY(POINT, 4326),
    last_seen_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ambulances_category_id ON ambulances(category_id);
CREATE INDEX idx_ambulances_station_facility_id ON ambulances(station_facility_id);
CREATE INDEX idx_ambulances_district_id ON ambulances(district_id);
CREATE INDEX idx_ambulances_status ON ambulances(status);
CREATE INDEX idx_ambulances_dispatch_readiness ON ambulances(dispatch_readiness);
CREATE INDEX idx_ambulances_location ON ambulances USING GIST(location);

CREATE TABLE ambulance_readiness_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ambulance_id UUID NOT NULL REFERENCES ambulances(id) ON DELETE CASCADE,
    mechanical_status TEXT NOT NULL CHECK (mechanical_status IN ('FUNCTIONAL', 'NEEDS_REPAIR', 'UNDER_MAINTENANCE', 'BREAKDOWN')),
    equipment_status TEXT NOT NULL CHECK (equipment_status IN ('FULLY_EQUIPPED', 'PARTIALLY_EQUIPPED', 'NOT_EQUIPPED')),
    fuel_status TEXT NOT NULL CHECK (fuel_status IN ('SUFFICIENT', 'LOW', 'CRITICAL')),
    oxygen_status TEXT NOT NULL CHECK (oxygen_status IN ('AVAILABLE', 'LOW', 'NOT_AVAILABLE')),
    tire_status TEXT NOT NULL CHECK (tire_status IN ('GOOD', 'WORN', 'BAD')),
    stretcher_status TEXT NOT NULL CHECK (stretcher_status IN ('AVAILABLE', 'DAMAGED', 'NOT_AVAILABLE')),
    communication_status TEXT NOT NULL CHECK (communication_status IN ('WORKING', 'INTERMITTENT', 'NOT_WORKING')),
    cleanliness_status TEXT NOT NULL CHECK (cleanliness_status IN ('CLEAN', 'DIRTY', 'NEEDS_SANITIZATION')),
    dispatch_readiness TEXT NOT NULL CHECK (dispatch_readiness IN ('DISPATCHABLE', 'RESTRICTED', 'NOT_DISPATCHABLE')),
    checked_by UUID REFERENCES users(id),
    checked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    notes TEXT
);

CREATE INDEX idx_ambulance_readiness_checks_ambulance_id ON ambulance_readiness_checks(ambulance_id, checked_at DESC);

CREATE TABLE ambulance_crew_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ambulance_id UUID NOT NULL REFERENCES ambulances(id) ON DELETE CASCADE,
    driver_user_id UUID REFERENCES users(id),
    medic_user_id UUID REFERENCES users(id),
    nurse_user_id UUID REFERENCES users(id),
    doctor_user_id UUID REFERENCES users(id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    released_at TIMESTAMPTZ,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    CONSTRAINT chk_ambulance_crew_times CHECK (released_at IS NULL OR released_at >= assigned_at)
);

CREATE INDEX idx_ambulance_crew_assignments_ambulance_id ON ambulance_crew_assignments(ambulance_id);
CREATE UNIQUE INDEX uq_ambulance_active_crew ON ambulance_crew_assignments(ambulance_id) WHERE active = TRUE;

CREATE TABLE ambulance_status_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ambulance_id UUID NOT NULL REFERENCES ambulances(id) ON DELETE CASCADE,
    previous_status TEXT,
    new_status TEXT NOT NULL,
    previous_dispatch_readiness TEXT,
    new_dispatch_readiness TEXT,
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    reason TEXT
);

CREATE INDEX idx_ambulance_status_logs_ambulance_id ON ambulance_status_logs(ambulance_id, changed_at DESC);

-- +goose Down
DROP TABLE IF EXISTS ambulance_status_logs;
DROP TABLE IF EXISTS ambulance_crew_assignments;
DROP TABLE IF EXISTS ambulance_readiness_checks;
DROP TABLE IF EXISTS ambulances;

