
-- ============================================
-- File: 000009_trips.sql
-- ============================================
-- +goose Up
CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dispatch_assignment_id UUID NOT NULL UNIQUE REFERENCES dispatch_assignments(id) ON DELETE CASCADE,
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    ambulance_id UUID REFERENCES ambulances(id),
    origin_lat DOUBLE PRECISION,
    origin_lon DOUBLE PRECISION,
    scene_lat DOUBLE PRECISION,
    scene_lon DOUBLE PRECISION,
    destination_facility_id UUID REFERENCES ref_facilities(id),
    destination_lat DOUBLE PRECISION,
    destination_lon DOUBLE PRECISION,
    odometer_start NUMERIC(12,2),
    odometer_end NUMERIC(12,2),
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    outcome TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_trips_odometer CHECK (odometer_end IS NULL OR odometer_start IS NULL OR odometer_end >= odometer_start),
    CONSTRAINT chk_trips_time CHECK (ended_at IS NULL OR started_at IS NULL OR ended_at >= started_at)
);

CREATE INDEX idx_trips_incident_id ON trips(incident_id);
CREATE INDEX idx_trips_ambulance_id ON trips(ambulance_id);
CREATE INDEX idx_trips_started_at ON trips(started_at DESC);

CREATE TABLE trip_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL CHECK (event_type IN ('ASSIGNED', 'ACCEPTED', 'DEPARTED', 'ARRIVED_SCENE', 'PATIENT_LOADED', 'ARRIVED_DESTINATION', 'COMPLETED', 'CANCELLED', 'GPS_PING', 'NOTE')),
    event_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    actor_user_id UUID REFERENCES users(id),
    notes TEXT
);

CREATE INDEX idx_trip_events_trip_id ON trip_events(trip_id, event_time DESC);

-- +goose Down
DROP TABLE IF EXISTS trip_events;
DROP TABLE IF EXISTS trips;

