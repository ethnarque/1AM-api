-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE countries (
    country_id SERIAL PRIMARY KEY,
    translation_key_id INTEGER NOT NULL REFERENCES keys(translation_key_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS countries;

-- +goose StatementEnd
