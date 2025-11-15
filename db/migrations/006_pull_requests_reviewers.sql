-- +goose Up
CREATE TABLE IF NOT EXISTS pull_request_reviewers (
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id     TEXT NOT NULL REFERENCES users(id),
    PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX IF NOT EXISTS idx_pull_request_reviewers_reviewer_id ON pull_request_reviewers(reviewer_id);

-- +goose Down
DROP INDEX IF EXISTS idx_pull_request_reviewers_reviewer_id;
DROP TABLE IF EXISTS pull_request_reviewers;