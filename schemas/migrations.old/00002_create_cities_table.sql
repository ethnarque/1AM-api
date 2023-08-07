-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE cities (
    city_id SERIAL PRIMARY KEY,
    translation_key_id INTEGER NOT NULL REFERENCES translation_key_id(translation_key_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS cities;

-- +goose StatementEnd
