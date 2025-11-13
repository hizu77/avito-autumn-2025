-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(name),
    is_active BOOLEAN NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_users_team_name;