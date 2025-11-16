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
) (model.ReassignedPullRequest, error) {
	pr, err := s.pullRequestStorage.GetPullRequestByID(ctx, id)
	if err != nil {
		return model.ReassignedPullRequest{}, errors.Wrap(err, "getting pull request")
	}

	if pr.Status == model.StatusMerged {
		return model.ReassignedPullRequest{}, model.ErrPullRequestIsMerged
	}

	if !slices.Contains(pr.ReviewersIDs, reviewerID) {
		return model.ReassignedPullRequest{}, model.ErrReviewerNotAssign
	}

	team, err := s.teamStorage.GetTeamByUserID(ctx, reviewerID)
	if err != nil {
		return model.ReassignedPullRequest{}, errors.Wrap(err, "getting team")
	}

	currentReviewers := make(map[string]struct{}, len(pr.ReviewersIDs))
	for _, id := range pr.ReviewersIDs {
		currentReviewers[id] = struct{}{}
	}

	validNewReviewers := collection.Filter(
		team.Members,
		func(user model.User) bool {
			if !user.IsActive || user.ID == reviewerID || user.ID == pr.AuthorID {
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
		return model.ReassignedPullRequest{}, errors.Wrap(err, "selecting new reviewer")
	}
	if len(selectedReviewers) == 0 {
		return model.ReassignedPullRequest{}, model.ErrNoCandidate
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
		updated, txErr := s.pullRequestStorage.UpdatePullRequestReviewers(ctx, pr)
		if txErr != nil {
			return errors.Wrap(txErr, "reassigning pull request")
		}

		updatedPr = updated

		return nil
	})
	if err != nil {
		return model.ReassignedPullRequest{}, errors.Wrap(err, "reassigning pull request in tx")
	}

	reassignedPullRequest := model.ReassignedPullRequest{
		ID:           updatedPr.ID,
		Name:         updatedPr.Name,
		AuthorID:     updatedPr.AuthorID,
		Status:       updatedPr.Status,
		ReviewersIDs: updatedPr.ReviewersIDs,
		CreatedAt:    updatedPr.CreatedAt,
		MergedAt:     updatedPr.MergedAt,
		ReassignedBy: newReviewerID,
	}

	return reassignedPullRequest, nil
}
