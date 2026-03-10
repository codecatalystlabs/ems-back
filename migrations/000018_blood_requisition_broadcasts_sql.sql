-- +goose Up
-- +goose StatementBegin
CREATE TABLE blood_requisition_broadcasts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blood_requisition_id UUID NOT NULL REFERENCES blood_requisitions(id) ON DELETE CASCADE,
    channel TEXT NOT NULL CHECK (channel IN ('SMS', 'APP', 'IN_APP', 'WHATSAPP', 'CALL')),
    recipient_site_id UUID REFERENCES blood_inventory_sites(id) ON DELETE SET NULL,
    recipient_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    recipient_phone TEXT,
    message_body TEXT NOT NULL,
    delivery_status TEXT NOT NULL DEFAULT 'PENDING'
        CHECK (delivery_status IN ('PENDING', 'SENT', 'DELIVERED', 'FAILED', 'ACKNOWLEDGED')),
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_blood_requisition_broadcasts_req_id
    ON blood_requisition_broadcasts(blood_requisition_id);
CREATE INDEX idx_blood_requisition_broadcasts_delivery_status
    ON blood_requisition_broadcasts(delivery_status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blood_requisition_broadcasts;
-- +goose StatementEnd
