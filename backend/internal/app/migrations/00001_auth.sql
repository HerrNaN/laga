-- +goose Up
CREATE TABLE users (
    id bytea PRIMARY KEY,
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE credentials (
    id bytea PRIMARY KEY,
    user_id bytea NOT NULL UNIQUE REFERENCES users(id),
    data jsonb NOT NULL
);

CREATE TABLE registration_attempts (
    token_hash bytea PRIMARY KEY,
    user_id bytea NOT NULL,
    name text NOT NULL,
    ceremony jsonb NOT NULL,
    expires_at timestamptz NOT NULL
);

CREATE TABLE sessions (
    token_hash bytea PRIMARY KEY,
    user_id bytea NOT NULL REFERENCES users(id),
    expires_at timestamptz NOT NULL
);

-- +goose Down
DROP TABLE sessions;
DROP TABLE registration_attempts;
DROP TABLE credentials;
DROP TABLE users;
