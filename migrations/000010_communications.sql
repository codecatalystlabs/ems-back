
-- ============================================
-- File: 000010_communications.sql
-- ============================================
-- +goose Up
CREATE TABLE inbound_sms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_phone TEXT NOT NULL,
    short_code TEXT,
    message_body TEXT NOT NULL,
    parsed_payload_json JSONB,
    parse_status TEXT NOT NULL DEFAULT 'PENDING' CHECK (parse_status IN ('PENDING', 'PARSED', 'FAILED')),
    linked_incident_id UUID REFERENCES incidents(id),
    provider_message_id TEXT,
    received_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_inbound_sms_sender_phone ON inbound_sms(sender_phone);
CREATE INDEX idx_inbound_sms_received_at ON inbound_sms(received_at DESC);

CREATE TABLE outbound_sms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_phone TEXT NOT NULL,
    message_body TEXT NOT NULL,
    purpose TEXT,
    provider TEXT,
    provider_message_id TEXT,
    status TEXT NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'SENT', 'DELIVERED', 'FAILED')),
    linked_incident_id UUID REFERENCES incidents(id),
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_outbound_sms_recipient_phone ON outbound_sms(recipient_phone);
CREATE INDEX idx_outbound_sms_status ON outbound_sms(status);

CREATE TABLE ussd_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id TEXT NOT NULL UNIQUE,
    phone_number TEXT NOT NULL,
    short_code TEXT,
    current_step TEXT,
    request_payload_json JSONB,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    linked_incident_id UUID REFERENCES incidents(id),
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ended_at TIMESTAMPTZ
);

CREATE INDEX idx_ussd_sessions_phone_number ON ussd_sessions(phone_number);
CREATE INDEX idx_ussd_sessions_started_at ON ussd_sessions(started_at DESC);

CREATE TABLE call_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number TEXT NOT NULL,
    call_direction TEXT NOT NULL CHECK (call_direction IN ('INBOUND', 'OUTBOUND')),
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    duration_seconds INT GENERATED ALWAYS AS (
        CASE WHEN ended_at IS NULL THEN NULL ELSE GREATEST(0, FLOOR(EXTRACT(EPOCH FROM (ended_at - started_at)))::INT) END
    ) STORED,
    agent_user_id UUID REFERENCES users(id),
    linked_incident_id UUID REFERENCES incidents(id),
    recording_url TEXT,
    outcome TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_call_logs_time CHECK (ended_at IS NULL OR ended_at >= started_at)
);

CREATE INDEX idx_call_logs_phone_number ON call_logs(phone_number);
CREATE INDEX idx_call_logs_started_at ON call_logs(started_at DESC);

-- +goose Down
DROP TABLE IF EXISTS call_logs;
DROP TABLE IF EXISTS ussd_sessions;
DROP TABLE IF EXISTS outbound_sms;
DROP TABLE IF EXISTS inbound_sms;
