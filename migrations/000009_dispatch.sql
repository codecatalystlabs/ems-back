-- ============================================
-- File: 000009_dispatch.sql
-- ============================================
-- +goose Up
CREATE TABLE dispatch_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    ambulance_id UUID REFERENCES ambulances(id),
    assigned_by_user_id UUID REFERENCES users(id),
    driver_user_id UUID REFERENCES users(id),
    lead_medic_user_id UUID REFERENCES users(id),
    team_snapshot_json JSONB,
    assignment_mode TEXT NOT NULL DEFAULT 'MANUAL' CHECK (assignment_mode IN ('MANUAL', 'ASSISTED', 'AUTO')),
    ranking_score NUMERIC(10,2),
    eta_minutes INT,
    status TEXT NOT NULL DEFAULT 'ASSIGNED' CHECK (status IN ('PROPOSED', 'ASSIGNED', 'ACCEPTED', 'DECLINED', 'DEPARTED', 'ARRIVED_SCENE', 'PATIENT_LOADED', 'ARRIVED_DESTINATION', 'COMPLETED', 'CANCELLED')),
    assigned_at TIMESTAMPTZ,
    accepted_at TIMESTAMPTZ,
    departed_at TIMESTAMPTZ,
    arrived_scene_at TIMESTAMPTZ,
    patient_loaded_at TIMESTAMPTZ,
    arrived_destination_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_dispatch_assignments_incident_id ON dispatch_assignments(incident_id);
CREATE INDEX idx_dispatch_assignments_ambulance_id ON dispatch_assignments(ambulance_id);
CREATE INDEX idx_dispatch_assignments_status ON dispatch_assignments(status);
CREATE INDEX idx_dispatch_assignments_assigned_by_user_id ON dispatch_assignments(assigned_by_user_id);
CREATE UNIQUE INDEX uq_open_dispatch_per_incident ON dispatch_assignments(incident_id)
WHERE status IN ('PROPOSED', 'ASSIGNED', 'ACCEPTED', 'DEPARTED', 'ARRIVED_SCENE', 'PATIENT_LOADED');

CREATE TABLE dispatch_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    ambulance_id UUID NOT NULL REFERENCES ambulances(id) ON DELETE CASCADE,
    driver_user_id UUID REFERENCES users(id),
    score NUMERIC(10,2) NOT NULL,
    eta_minutes INT,
    rule_summary TEXT,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    selected BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_dispatch_recommendations_incident_id ON dispatch_recommendations(incident_id, generated_at DESC);

ALTER TABLE user_availability
    ADD CONSTRAINT fk_user_availability_current_incident
    FOREIGN KEY (current_incident_id) REFERENCES incidents(id);

ALTER TABLE user_availability
    ADD CONSTRAINT fk_user_availability_current_dispatch_assignment
    FOREIGN KEY (current_dispatch_assignment_id) REFERENCES dispatch_assignments(id);

-- +goose Down
ALTER TABLE user_availability DROP CONSTRAINT IF EXISTS fk_user_availability_current_dispatch_assignment;
ALTER TABLE user_availability DROP CONSTRAINT IF EXISTS fk_user_availability_current_incident;
DROP TABLE IF EXISTS dispatch_recommendations;
DROP TABLE IF EXISTS dispatch_assignments;
