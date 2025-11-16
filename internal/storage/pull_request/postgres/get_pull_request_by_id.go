package pullrequest

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/pull_request/dbmodel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (s *Storage) GetPullRequestByID(ctx context.Context, id string) (model.PullRequest, error) {
	sql, args, err := squirrel.
		Expr(`
        	SELECT
            	pr.id         AS pr_id,
            	pr.name       AS pr_name,
            	pr.author_id  AS author_id,
            	s.name        AS status,
            	pr.created_at AS created_at,
            	pr.merged_at  AS merged_at,
				COALESCE(
  					array_agg(r.reviewer_id ORDER BY r.reviewer_id)
    				FILTER (WHERE r.reviewer_id IS NOT NULL),
  					'{}'
				) AS reviewer_ids
			FROM pull_requests pr
        	JOIN pull_request_statuses s ON s.id = pr.status_id
        	LEFT JOIN pull_request_reviewers r ON r.pull_request_id = pr.id
        	WHERE pr.id = $1
        	GROUP BY
            	pr_id,
            	pr_name,
            	author_id,
            	status,
            	created_at,
            	merged_at
        `, id).
		ToSql()
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "building sql")
	}

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sql, args...)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "querying sql")
	}

	dbPR, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbmodel.PullRequest])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.PullRequest{}, model.ErrPullRequestDoesNotExist
	}
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "collecting rows")
	}

	pr, err := mapDBPullRequestToDomainPullRequest(dbPR)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "mapping pull request")
	}

	return pr, nil
}
