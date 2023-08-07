-- +goose Up
-- +goose StatementBegin
-- Create the locales table
CREATE TABLE locales (
    locale_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE CHECK (name = LOWER(name))
);

-- Create the keys table
CREATE TABLE translation_keys (
    translation_key_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    namespace VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create the translations table
CREATE TABLE translations (
    translation_id SERIAL PRIMARY KEY,
    translation_key_id INTEGER NOT NULL REFERENCES keys(translation_key_id),
    locale_id INTEGER NOT NULL REFERENCES locales(locale_id),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS translations;
DROP TABLE IF EXISTS keys;
DROP TABLE IF EXISTS locales;

-- +goose StatementEnd
