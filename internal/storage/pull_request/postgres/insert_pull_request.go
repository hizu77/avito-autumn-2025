package pullrequest

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/common/constraint"
	"github.com/pkg/errors"
)

func (s *Storage) InsertPullRequest(
	ctx context.Context,
	request model.PullRequest,
) (model.PullRequest, error) {
	createdAt := time.Now().UTC()
	request.CreatedAt = &createdAt

	sql, args, err := squirrel.
		Expr(`
			INSERT INTO pull_requests (
				id,
				name,
				author_id,
				status_id,
				created_at,
				merged_at
			)
			VALUES (
				$1,
				$2,
				$3,
				(SELECT id FROM pull_request_statuses WHERE name = $4),
				$5,
				$6
			)
		`,
			request.ID,
			request.Name,
			request.AuthorID,
			request.Status,
			createdAt,
			request.MergedAt,
		).ToSql()
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "building sql")
	}

	tx := s.getter.DefaultTrOrDB(ctx, s.pool)
	_, err = tx.Exec(ctx, sql, args...)
	if constraint.IsUniqueViolation(err) {
		return model.PullRequest{}, model.ErrPullRequestAlreadyExists
	}
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "executing sql")
	}

	if len(request.ReviewersIDs) == 0 {
		return request, nil
	}

	sql, args, err = squirrel.
		Expr(`
			INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
			SELECT $1, unnest($2::text[])
			ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING`,
			request.ID, request.ReviewersIDs).
		ToSql()
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "executing sql")
	}

	return request, nil
}
