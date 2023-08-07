-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    activated BOOL NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
