CREATE TABLE IF NOT EXISTS pull_requests (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    author_id  TEXT NOT NULL REFERENCES users(id),
    status_id  INTEGER NOT NULL REFERENCES pull_request_statuses(id),
    created_at TIMESTAMPTZ NOT NULL,
    merged_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id);

CREATE INDEX IF NOT EXISTS idx_pull_requests_status_id ON pull_requests(status_id);