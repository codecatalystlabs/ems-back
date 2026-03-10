-- +goose Up
-- +goose StatementBegin
CREATE TABLE blood_requisition_status_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blood_requisition_id UUID NOT NULL REFERENCES blood_requisitions(id) ON DELETE CASCADE,
    previous_status TEXT,
    new_status TEXT NOT NULL,
    actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_blood_requisition_status_logs_req_id
    ON blood_requisition_status_logs(blood_requisition_id, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blood_requisition_status_logs;
-- +goose StatementEnd
