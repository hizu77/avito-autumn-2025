-- +goose Up
CREATE TABLE IF NOT EXISTS pull_request_statuses (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO pull_request_statuses (id, name)
VALUES
    (1, 'open'),
    (2, 'merged');

-- +goose Down
DELETE FROM pull_request_statuses WHERE name IN ('open', 'merged');

DROP TABLE IF EXISTS pull_request_statuses;