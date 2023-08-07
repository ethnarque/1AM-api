-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS tracks (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP(0) WITH TIME ZONE,
    title VARCHAR(100) NOT NULL,
    duration INTEGER NOT NULL,
    genres TEXT[] NOT NULL,
    albums TEXT[] NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4 ()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS tracks;

-- +goose StatementEnd
