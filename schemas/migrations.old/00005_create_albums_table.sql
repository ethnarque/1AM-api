-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE albums (
    album_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    artist_id INTEGER NOT NULL REFERENCES artists(artist_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS albums;

-- +goose StatementEnd
