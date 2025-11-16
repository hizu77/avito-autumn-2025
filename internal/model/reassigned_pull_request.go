package model

import "time"

type ReassignedPullRequest struct {
	ID           string
	Name         string
	AuthorID     string
	Status       Status
	ReviewersIDs []string
	ReassignedBy string

	CreatedAt *time.Time
	MergedAt  *time.Time
}
