package dbmodel

import "time"

type PullRequest struct {
	ID          string   `db:"pr_id"`
	Name        string   `db:"pr_name"`
	AuthorID    string   `db:"author_id"`
	Status      string   `db:"status"`
	ReviewerIDs []string `db:"reviewer_ids"`

	CreatedAt time.Time  `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
}
