package pullrequest

import (
	"context"
	"slices"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/pkg/errors"
)

const (
	reassignReviewersCount = 1
)

func (s *Service) ReassignPullRequest(
	ctx context.Context,
	id string,
	reviewerID string,
) (model.PullRequest, error) {
	pr, err := s.pullRequestStorage.GetPullRequestByID(ctx, id)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "getting pull request")
	}

	if pr.Status == model.StatusMerged {
		return model.PullRequest{}, model.ErrPullRequestIsMerged
	}

	if !slices.Contains(pr.ReviewersIDs, reviewerID) {
		return model.PullRequest{}, model.ErrReviewerNotAssign
	}

	team, err := s.teamStorage.GetTeamByUserID(ctx, reviewerID)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "getting team")
	}

	currentReviewers := make(map[string]struct{}, len(pr.ReviewersIDs))
	for _, id := range pr.ReviewersIDs {
		currentReviewers[id] = struct{}{}
	}

	validNewReviewers := collection.Filter(
		team.Members,
		func(user model.User) bool {
			if !(user.IsActive && user.ID != reviewerID) {
				return false
			}

			if _, exists := currentReviewers[user.ID]; exists {
				return false
			}

			return true
		},
	)
	validNewReviewersIDs := collection.Map(validNewReviewers, model.User.GetID)
	selectedReviewers, err := s.getRandomReviewers(
		validNewReviewersIDs,
		reassignReviewersCount,
	)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "selecting new reviewer")
	}
	if len(selectedReviewers) == 0 {
		return model.PullRequest{}, model.ErrNoCandidate
	}

	newReviewerID := selectedReviewers[0]

	for i, id := range pr.ReviewersIDs {
		if id == reviewerID {
			pr.ReviewersIDs[i] = newReviewerID
			break
		}
	}

	var updatedPr model.PullRequest
	err = s.trManager.Do(ctx, func(ctx context.Context) error {
		updated, err := s.pullRequestStorage.UpdatePullRequestReviewers(ctx, pr)
		if err != nil {
			return errors.Wrap(err, "reassigning pull request")
		}

		updatedPr = updated

		return nil
	})
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "reassigning pull request in tx")
	}

	return updatedPr, nil
}
