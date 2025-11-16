package model

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrPullRequestAlreadyExists = errors.New("pull request already exists")
	ErrPullRequestDoesNotExist  = errors.New("pull request does not exist")
	ErrPullRequestIsMerged      = errors.New("pull request is merged")
	ErrReviewerNotAssign        = errors.New("not assigned reviewer")
	ErrNoCandidate              = errors.New("no candidate to reassign")
)

type PullRequest struct {
	ID           string
	Name         string
	AuthorID     string
	Status       Status
	ReviewersIDs []string

	CreatedAt *time.Time
	MergedAt  *time.Time
}
