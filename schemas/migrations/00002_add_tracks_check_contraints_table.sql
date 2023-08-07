-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE tracks
    ADD CONSTRAINT tracks_duration_check
    CHECK (duration >= 0);

ALTER TABLE tracks 
    ADD CONSTRAINT genres_length_check 
    CHECK (array_length(genres, 1) BETWEEN 1 AND 5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE tracks
    DROP CONSTRAINT IF EXISTS tracks_duration_check;

ALTER TABLE tracks
    DROP CONSTRAINT IF EXISTS tracks_genres_check;
-- +goose StatementEnd
