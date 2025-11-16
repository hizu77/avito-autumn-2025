CREATE TABLE IF NOT EXISTS pull_request_statuses (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO pull_request_statuses (id, name)
VALUES
    (1, 'OPEN'),
    (2, 'MERGED');