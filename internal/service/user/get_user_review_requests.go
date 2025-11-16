package user

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
)

func (s *Service) GetUserReviewRequests(ctx context.Context, id string) ([]model.PullRequest, error) {
	return s.pullRequestStorage.GetPullRequestsByReviewer(ctx, id)
}
