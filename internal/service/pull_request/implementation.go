package pullrequest

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/hizu77/avito-autumn-2025/internal/model"
)

type (
	teamStorage interface {
		GetTeamByUserID(ctx context.Context, userID string) (model.Team, error)
	}

	pullRequestStorage interface {
		GetPullRequestByID(ctx context.Context, id string) (model.PullRequest, error)
		InsertPullRequest(ctx context.Context, request model.PullRequest) (model.PullRequest, error)
		MergePullRequest(ctx context.Context, req model.PullRequest) (model.PullRequest, error)
		ReassignReviewer(
			ctx context.Context,
			req model.PullRequest,
			oldReviewer string,
			newReviewer string,
		) (model.PullRequest, error)
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
