-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE tracks (
    track_id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    duration TIME NOT NULL,
    album_id INTEGER NOT NULL REFERENCES albums(album_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS tracks;

-- +goose StatementEnd
