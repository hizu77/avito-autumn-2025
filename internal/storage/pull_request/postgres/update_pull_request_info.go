package pullrequest

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
)

func (s *Storage) UpdatePullRequestInfo(
	ctx context.Context,
	req model.PullRequest,
) (model.PullRequest, error) {
	sql, args, err := squirrel.
		Expr(`
        UPDATE pull_requests
        SET
            name       = $1,
            author_id  = $2,
            status_id  = (SELECT id FROM pull_request_statuses WHERE name = $3),
            created_at = $4,
            merged_at  = $5
        WHERE id = $6
    	`,
			req.Name,
			req.AuthorID,
			req.Status,
			req.CreatedAt,
			req.MergedAt,
			req.ID,
		).
		ToSql()
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "building sql")
	}

	tag, err := s.getter.DefaultTrOrDB(ctx, s.pool).Exec(ctx, sql, args...)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "executing sql")
	}
	if tag.RowsAffected() == 0 {
		return model.PullRequest{}, model.ErrPullRequestDoesNotExist
	}

	return req, nil
}
