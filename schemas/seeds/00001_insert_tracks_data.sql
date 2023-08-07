-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO tracks (title, duration, genres, albums) VALUES
    ('Yugoslavskiy Groove', '241', '{"electronics"}', '{"Yugoslavskiy Groove"}'),
    ('Diarabi', '301', '{"africa", "malian"}', '{"Ali"}');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DELETE FROM tracks;
-- +goose StatementEnd
