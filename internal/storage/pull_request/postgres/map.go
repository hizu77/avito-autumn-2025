package pullrequest

import (
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/pull_request/dbmodel"
	"github.com/pkg/errors"
)

func mapDBPullRequestToDomainPullRequest(pr dbmodel.PullRequest) (model.PullRequest, error) {
	mappedStatus, err := model.ParseStatus(pr.Status)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "mapping status")
	}

	return model.PullRequest{
		ID:           pr.ID,
		Name:         pr.Name,
		AuthorID:     pr.AuthorID,
		Status:       mappedStatus,
		ReviewersIDs: pr.ReviewerIDs,
		CreatedAt:    &pr.CreatedAt,
		MergedAt:     pr.MergedAt,
	}, nil
}
