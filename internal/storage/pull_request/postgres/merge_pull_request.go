package pullrequest

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
)

func (s *Storage) MergePullRequest(
	ctx context.Context,
	req model.PullRequest,
) (model.PullRequest, error) {
	mergedAt := time.Now().UTC()
	sql, args, err := squirrel.
		Expr(`
			UPDATE pull_requests
			SET status_id = (SELECT id FROM pull_request_statuses WHERE name = $1),
			    merged_at = $2
			WHERE id = $3
		`,
			model.StatusMerged,
			mergedAt,
			req.ID,
		).
		ToSql()
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "buidling sql")
	}

	tag, err := s.getter.DefaultTrOrDB(ctx, s.pool).Exec(ctx, sql, args...)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "executing sql")
	}
	if tag.RowsAffected() == 0 {
		return model.PullRequest{}, model.ErrPullRequestDoesNotExist
	}

	req.MergedAt = &mergedAt
	req.Status = model.StatusMerged

	return req, nil
}
