-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX IF NOT EXISTS tracks_title_idx ON tracks USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS tracks_genres_idx ON tracks USING GIN(genres);
CREATE INDEX IF NOT EXISTS tracks_albums_idx ON tracks USING GIN(albums);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS tracks_title_idx;
DROP INDEX IF EXISTS tracks_genres_idx;
DROP INDEX IF EXISTS tracks_albums_idx;

-- +goose StatementEnd
