package pullrequest

import (
	"context"
	"time"

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

	now := time.Now().UTC()
	pr.Status = model.StatusMerged
	pr.MergedAt = &now

	updated, err := s.pullRequestStorage.UpdatePullRequestInfo(ctx, pr)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "updating pull request info")
	}

	return updated, nil
}
