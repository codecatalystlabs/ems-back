-- ============================================
-- File: 000011_notifications.sql
-- ============================================
-- +goose Up
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL,
    recipient_user_id UUID REFERENCES users(id),
    recipient_phone TEXT,
    recipient_email TEXT,
    title TEXT,
    body TEXT NOT NULL,
    channel TEXT NOT NULL CHECK (channel IN ('SMS', 'EMAIL', 'PUSH', 'IN_APP')),
    linked_entity_type TEXT,
    linked_entity_id UUID,
    status TEXT NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'SENT', 'DELIVERED', 'READ', 'FAILED')),
    attempts INT NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    sent_at TIMESTAMPTZ,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_recipient_user_id ON notifications(recipient_user_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_channel ON notifications(channel);

-- +goose Down
DROP TABLE IF EXISTS notifications;
