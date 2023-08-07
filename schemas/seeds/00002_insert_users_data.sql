-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- the password is 'pa55word'
INSERT INTO users (name, email, password_hash, activated) VALUES
    ('Jane Doe', 'jane@example.com', '$2y$10$F1cCrYdGP7sVjB9ONqyOjOHyC22n.KY3x1y1aEAONjIV2XIPhmRgu', 'true'),
    ('Alice Doe', 'alice@example.com', '$2y$10$bZsbyO23GoEljxUw4pB.S.YCUjE17eW3.EeXo5tjS2uS7/3vqFNvK', 'true'),
    ('Michael Scott', 'michael@example.com', '$2y$10$/huQSzH8fnWacywQG9VhFubLQ1nctvy8yNHWNtQmJ3c3tGxqHhzD.', 'false');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DELETE FROM users;

-- +goose StatementEnd
