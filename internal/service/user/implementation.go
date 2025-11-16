package user

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
)

//go:generate mockgen -source=implementation.go -destination=../../mock/user/storage.go -package=mock -mock_names userStorage=UserStorage,pullRequestStorage=PullRequestStorage
type (
	userStorage interface {
		UpdateActivity(ctx context.Context, id string, activity bool) (model.User, error)
	}

	pullRequestStorage interface {
		GetPullRequestsByReviewer(ctx context.Context, id string) ([]model.PullRequest, error)
	}
)

type Service struct {
	userStorage        userStorage
	pullRequestStorage pullRequestStorage
}

func New(
	userStorage userStorage,
	pullRequestStorage pullRequestStorage,
) *Service {
	return &Service{
		userStorage:        userStorage,
		pullRequestStorage: pullRequestStorage,
	}
}
