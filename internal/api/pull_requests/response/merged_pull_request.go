package response

import (
	"time"

	"github.com/hizu77/avito-autumn-2025/internal/model"
)

type MergedPullRequest struct {
	ID        string       `json:"pull_request_id"`
	Name      string       `json:"pull_request_name"`
	AuthorID  string       `json:"author_id"`
	Status    model.Status `json:"status"`
	Reviewers []string     `json:"assigned_reviewers"`
	MergedAt  time.Time    `json:"merged_at"`
}
