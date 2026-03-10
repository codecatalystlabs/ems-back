-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';CREATE TABLE blood_transport_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blood_requisition_id UUID NOT NULL REFERENCES blood_requisitions(id) ON DELETE CASCADE,
    blood_requisition_offer_id UUID REFERENCES blood_requisition_offers(id) ON DELETE SET NULL,
    vehicle_type TEXT NOT NULL CHECK (vehicle_type IN ('AMBULANCE', 'PICKUP', 'MOTORCYCLE', 'OTHER')),
    ambulance_id UUID REFERENCES ambulances(id) ON DELETE SET NULL,
    dispatch_assignment_id UUID REFERENCES dispatch_assignments(id) ON DELETE SET NULL,
    assigned_driver_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    pickup_site_id UUID REFERENCES blood_inventory_sites(id) ON DELETE SET NULL,
    destination_facility_id UUID REFERENCES ref_facilities(id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'ASSIGNED'
        CHECK (status IN ('ASSIGNED', 'ACCEPTED', 'ENROUTE_PICKUP', 'COLLECTED', 'ENROUTE_DESTINATION', 'DELIVERED', 'CANCELLED')),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    collected_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    notes TEXT
);

CREATE INDEX idx_blood_transport_assignments_req_id
    ON blood_transport_assignments(blood_requisition_id);
CREATE INDEX idx_blood_transport_assignments_offer_id
    ON blood_transport_assignments(blood_requisition_offer_id);
CREATE INDEX idx_blood_transport_assignments_status
    ON blood_transport_assignments(status);
CREATE INDEX idx_blood_transport_assignments_ambulance_id
    ON blood_transport_assignments(ambulance_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blood_transport_assignments;
-- +goose StatementEnd
