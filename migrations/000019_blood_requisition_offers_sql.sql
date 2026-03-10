-- +goose Up
-- +goose StatementBegin
CREATE TABLE blood_requisition_offers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blood_requisition_id UUID NOT NULL REFERENCES blood_requisitions(id) ON DELETE CASCADE,
    inventory_site_id UUID NOT NULL REFERENCES blood_inventory_sites(id) ON DELETE CASCADE,
    blood_product_id UUID NOT NULL REFERENCES blood_products(id),
    blood_group_id UUID NOT NULL REFERENCES blood_groups(id),
    units_offered INT NOT NULL CHECK (units_offered > 0),
    reserved_until TIMESTAMPTZ,
    notes TEXT,
    contact_person_name TEXT,
    contact_phone TEXT,
    offered_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'OFFERED'
        CHECK (status IN ('OFFERED', 'ACCEPTED', 'DECLINED', 'EXPIRED', 'FULFILLED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_blood_requisition_offers_req_id
    ON blood_requisition_offers(blood_requisition_id);
CREATE INDEX idx_blood_requisition_offers_site_id
    ON blood_requisition_offers(inventory_site_id);
CREATE INDEX idx_blood_requisition_offers_status
    ON blood_requisition_offers(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blood_requisition_offers;
-- +goose StatementEnd
