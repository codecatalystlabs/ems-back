-- +goose Up
-- +goose StatementBegin
CREATE TABLE blood_requisitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID REFERENCES incidents(id) ON DELETE SET NULL,
    requesting_facility_id UUID REFERENCES ref_facilities(id) ON DELETE SET NULL,
    patient_name TEXT,
    patient_identifier TEXT,
    clinical_summary TEXT NOT NULL,
    diagnosis TEXT,
    indication TEXT,
    parity_summary TEXT,
    blood_group_id UUID NOT NULL REFERENCES blood_groups(id),
    blood_product_id UUID NOT NULL REFERENCES blood_products(id),
    units_requested INT NOT NULL CHECK (units_requested > 0),
    urgency_level TEXT NOT NULL DEFAULT 'EMERGENCY'
        CHECK (urgency_level IN ('EMERGENCY', 'URGENT', 'ROUTINE')),
    status TEXT NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'BROADCASTING', 'MATCHED', 'PICKUP_ASSIGNED', 'COLLECTED', 'DELIVERED', 'CANCELLED', 'EXPIRED')),
    reporter_phone TEXT,
    destination_facility_id UUID REFERENCES ref_facilities(id) ON DELETE SET NULL,
    destination_lat DOUBLE PRECISION,
    destination_lon DOUBLE PRECISION,
    destination_location GEOGRAPHY(POINT, 4326),
    requested_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_blood_requisitions_incident_id ON blood_requisitions(incident_id);
CREATE INDEX idx_blood_requisitions_requesting_facility_id ON blood_requisitions(requesting_facility_id);
CREATE INDEX idx_blood_requisitions_destination_facility_id ON blood_requisitions(destination_facility_id);
CREATE INDEX idx_blood_requisitions_status ON blood_requisitions(status);
CREATE INDEX idx_blood_requisitions_created_at ON blood_requisitions(created_at DESC);
CREATE INDEX idx_blood_requisitions_destination_location ON blood_requisitions USING GIST(destination_location);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blood_requisitions;
-- +goose StatementEnd
