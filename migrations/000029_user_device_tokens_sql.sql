-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_device_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id TEXT,
    device_name TEXT,
    platform TEXT NOT NULL CHECK (platform IN ('ANDROID', 'IOS', 'WEB')),
    push_token TEXT NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_device_tokens_user_id ON user_device_tokens(user_id);
CREATE INDEX idx_user_device_tokens_platform ON user_device_tokens(platform);
CREATE INDEX idx_user_device_tokens_is_active ON user_device_tokens(is_active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_device_tokens;
-- +goose StatementEnd
