package pullrequest

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
)

func (s *Storage) ReassignReviewer(
	ctx context.Context,
	req model.PullRequest,
	oldReviewer string,
	newReviewer string,
) (model.PullRequest, error) {
	sql, args, err := squirrel.
		Delete(pullRequestReviewersTable).
		Where(squirrel.Eq{
			columnPullRequestID: req.ID,
			columnReviewerID:    oldReviewer,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "building sql")
	}

	tx := s.getter.DefaultTrOrDB(ctx, s.pool)
	tag, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "executing sql")
	}
	if tag.RowsAffected() == 0 {
		return model.PullRequest{}, model.ErrReviewerNotAssign
	}

	sql, args, err = squirrel.
		Insert(pullRequestReviewersTable).
		Columns(columnPullRequestID, columnReviewerID).
		Values(req.ID, newReviewer).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "building sql")
	}

	tag, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return model.PullRequest{}, errors.Wrap(err, "executing sql")
	}

	updated := req
	for i, id := range updated.ReviewersIDs {
		if id == oldReviewer {
			updated.ReviewersIDs[i] = newReviewer
			return updated, nil
		}
	}

	return updated, nil
}
