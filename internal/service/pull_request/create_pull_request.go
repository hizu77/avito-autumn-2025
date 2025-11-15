package pullrequest

import (
	"context"
	"math/rand"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/pkg/errors"
)

const (
	maxReviewerCntWithAuthor = 2
)

func (s *Service) CreatePullRequest(ctx context.Context, request model.PullRequest) (model.PullRequest, error) {
	var createdPullRequest model.PullRequest
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		team, err := s.teamStorage.GetTeamByUserID(ctx, request.AuthorID)
		if err != nil {
			return errors.Wrap(err, "getting team by user ID")
		}

		activeTeammates := collection.Filter(
			team.Members,
			func(user model.User) bool {
				return user.IsActive && user.ID != request.AuthorID
			},
		)
		reviewers := make(map[string]struct{})
		reviewersCnt := min(maxReviewerCntWithAuthor, len(activeTeammates))
		for len(reviewers) < reviewersCnt {
			randomIdx := rand.Intn(len(activeTeammates))
			reviewer := activeTeammates[randomIdx].ID

			if _, ok := reviewers[reviewer]; !ok {
				reviewers[reviewer] = struct{}{}
			}
		}

		assignedReviewers := collection.Keys(reviewers)
		pullRequest := model.PullRequest{
			ID:           request.ID,
			Name:         request.Name,
			AuthorID:     request.AuthorID,
			Status:       model.StatusOpen,
			ReviewersIDs: assignedReviewers,
		}

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
