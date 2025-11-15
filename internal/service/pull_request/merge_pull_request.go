package pullrequest

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
)

func (s *Service) MergePullRequest(ctx context.Context, id string) (model.PullRequest, error) {
	pr, err := s.pullRequestStorage.GetPullRequestByID(ctx, id)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "getting pull request")
	}

	if pr.Status == model.StatusMerged {
		return pr, nil
	}

	merged, err := s.pullRequestStorage.MergePullRequest(ctx, pr)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "merging pull request")
	}

	return merged, nil
}
