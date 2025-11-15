package pullrequest

import (
	"context"
	"math/rand"
	"slices"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/pkg/errors"
)

func (s *Service) ReassignPullRequest(
	ctx context.Context,
	id string,
	reviewerID string,
) (model.PullRequest, error) {
	var reassignedPullRequest model.PullRequest
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		pr, err := s.pullRequestStorage.GetPullRequestByID(ctx, id)
		if err != nil {
			return errors.Wrap(err, "getting pull request")
		}

		if !slices.Contains(pr.ReviewersIDs, reviewerID) {
			return model.ErrReviewerNotAssign
		}

		team, err := s.teamStorage.GetTeamByUserID(ctx, reviewerID)
		if err != nil {
			return errors.Wrap(err, "getting team")
		}

		activeTeammates := collection.Filter(
			team.Members,
			func(user model.User) bool {
				return user.IsActive && user.ID != reviewerID
			},
		)
		if len(activeTeammates) == 0 {
			return model.ErrNoCandidate
		}

		currentReviewers := make(map[string]struct{}, len(pr.ReviewersIDs))
		for _, id := range pr.ReviewersIDs {
			currentReviewers[id] = struct{}{}
		}

		candidates := collection.Filter(
			activeTeammates,
			func(u model.User) bool {
				_, alreadyReviewer := currentReviewers[u.ID]
				return !alreadyReviewer
			},
		)
		if len(candidates) == 0 {
			return model.ErrNoCandidate
		}

		newIdx := rand.Intn(len(candidates))
		newReviewerID := candidates[newIdx].ID

		updated, err := s.pullRequestStorage.ReassignReviewer(ctx, pr, reviewerID, newReviewerID)
		if err != nil {
			return errors.Wrap(err, "reassigning pull request")
		}

		reassignedPullRequest = updated

		return nil
	})
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "reassigning pull request in tx")
	}

	return reassignedPullRequest, nil
}
