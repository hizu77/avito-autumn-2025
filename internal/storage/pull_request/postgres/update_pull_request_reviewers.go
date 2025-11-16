package pullrequest

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
)

func (s *Storage) UpdatePullRequestReviewers(
	ctx context.Context,
	req model.PullRequest,
) (model.PullRequest, error) {
	sql, args, err := squirrel.
		Delete(pullRequestReviewersTable).
		Where(squirrel.Eq{columnPullRequestID: req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "building delete sql")
	}

	tx := s.getter.DefaultTrOrDB(ctx, s.pool)

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "executing delete sql")
	}

	if len(req.ReviewersIDs) == 0 {
		return req, nil
	}

	sql, args, err = squirrel.
		Expr(`
			INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
			SELECT $1, unnest($2::text[])
			ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING
		`,
			req.ID,
			req.ReviewersIDs,
		).
		ToSql()
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "building insert sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "executing insert sql")
	}

	return req, nil
}
