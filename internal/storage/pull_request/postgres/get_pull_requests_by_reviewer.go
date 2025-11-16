package pullrequest

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/pull_request/dbmodel"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (s *Storage) GetPullRequestsByReviewer(ctx context.Context, id string) ([]model.PullRequest, error) {
	sql, args, err := squirrel.
		Expr(`
        	SELECT
            	pr.id 															  AS pr_id,
            	pr.name 														  AS pr_name,
            	pr.author_id 													  AS author_id,
            	s.name 															  AS status,
            	pr.created_at 													  AS created_at,
            	pr.merged_at 													  AS merged_at,
            	COALESCE(array_agg(r2.reviewer_id ORDER BY r2.reviewer_id), '{}') AS reviewer_ids
        	FROM pull_requests pr 
			JOIN pull_request_reviewers prr ON prr.pull_request_id = pr.id
        	JOIN pull_request_statuses s ON s.id = pr.status_id
        	LEFT JOIN pull_request_reviewers r2 ON r2.pull_request_id = pr.id
        	WHERE prr.reviewer_id = $1
        	GROUP BY
            	pr_id,
            	pr_name,
            	author_id,
            	status,
            	created_at,
            	merged_at`, id).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying sql")
	}

	fetched, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.PullRequest])
	if err != nil {
		return nil, errors.Wrap(err, "collecting rows")
	}
	if len(fetched) == 0 {
		return []model.PullRequest{}, nil
	}

	mappedPullRequests, err := collection.MapWithError(fetched, mapDBPullRequestToDomainPullRequest)
	if err != nil {
		return nil, errors.Wrap(err, "mapping pull requests")
	}

	return mappedPullRequests, nil
}
