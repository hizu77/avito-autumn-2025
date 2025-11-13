-- +goose Up
CREATE TABLE IF NOT EXISTS admins (
    id TEXT PRIMARY KEY,
    password TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS admins;