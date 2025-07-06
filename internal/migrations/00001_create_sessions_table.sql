-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS refresh_sessions (
    user_id INT,
    token_hash VARCHAR(255),
    user_agent VARCHAR(255),
    ip_address VARCHAR(255),
    expiry_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '1 day',
    is_revoked BOOLEAN
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_sessions;
-- +goose StatementEnd
