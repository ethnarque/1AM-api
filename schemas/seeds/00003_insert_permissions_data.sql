-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO permissions (code) VALUES
('tracks:read'), ('tracks:write');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';


DELETE FROM permissions;
-- +goose StatementEnd
