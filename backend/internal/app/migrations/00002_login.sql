-- +goose Up
CREATE TABLE login_attempts (
    token_hash bytea PRIMARY KEY,
    ceremony jsonb NOT NULL,
    expires_at timestamptz NOT NULL
);

-- +goose Down
DROP TABLE login_attempts;
