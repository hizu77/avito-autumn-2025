package pullrequest

import (
	"context"
	"time"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/pkg/errors"
)

const (
	maxCreateReviewersCount = 2
)

func (s *Service) CreatePullRequest(ctx context.Context, request model.PullRequest) (model.PullRequest, error) {
	team, err := s.teamStorage.GetTeamByUserID(ctx, request.AuthorID)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "getting team by user ID")
	}

	activeTeammates := collection.Filter(
		team.Members,
		func(user model.User) bool {
			return user.IsActive && user.ID != request.AuthorID
		},
	)
	activeTeammatesIDs := collection.Map(activeTeammates, model.User.GetID)

	reviewers, err := s.getRandomReviewers(activeTeammatesIDs, maxCreateReviewersCount)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "getting random reviewers")
	}

	createdAt := time.Now().UTC()
	pullRequest := model.PullRequest{
		ID:           request.ID,
		Name:         request.Name,
		AuthorID:     request.AuthorID,
		Status:       model.StatusOpen,
		ReviewersIDs: reviewers,
		CreatedAt:    &createdAt,
	}

	var createdPullRequest model.PullRequest
	err = s.trManager.Do(ctx, func(ctx context.Context) error {
		inserted, err := s.pullRequestStorage.InsertPullRequest(ctx, pullRequest)
		if err != nil {
			return errors.Wrap(err, "insert pull request")
		}

		createdPullRequest = inserted

		return nil
	})
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "creating pull request in tx")
	}

	return createdPullRequest, nil
}
