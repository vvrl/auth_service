-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS refresh_sessions (
    id SERIAL PRIMARY KEY, 
    user_id UUID NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    ip_address VARCHAR(255) NOT NULL,
    expiry_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '1 day',
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_refresh_sessions_user_id ON refresh_sessions(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_sessions;
-- +goose StatementEnd
