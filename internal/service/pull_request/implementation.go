package pullrequest

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/hizu77/avito-autumn-2025/internal/model"
)

//go:generate mockgen -source=implementation.go -destination=../../mock/pull_request/storage.go -package=mock -mock_names teamStorage=TeamStorage,pullRequestStorage=PullRequestStorage
type (
	teamStorage interface {
		GetTeamByUserID(ctx context.Context, userID string) (model.Team, error)
	}

	pullRequestStorage interface {
		GetPullRequestByID(ctx context.Context, id string) (model.PullRequest, error)
		InsertPullRequest(ctx context.Context, request model.PullRequest) (model.PullRequest, error)
		UpdatePullRequestInfo(ctx context.Context, req model.PullRequest) (model.PullRequest, error)
		UpdatePullRequestReviewers(ctx context.Context, req model.PullRequest) (model.PullRequest, error)
	}
)

type Service struct {
	teamStorage        teamStorage
	pullRequestStorage pullRequestStorage

	trManager trm.Manager
}

func New(
	teamStorage teamStorage,
	pullRequestStorage pullRequestStorage,
	trManager trm.Manager,
) *Service {
	return &Service{
		teamStorage:        teamStorage,
		pullRequestStorage: pullRequestStorage,
		trManager:          trManager,
	}
}
