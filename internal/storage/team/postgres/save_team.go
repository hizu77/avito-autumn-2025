package team

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/common/constraint"
	"github.com/pkg/errors"
)

func (s *Storage) SaveTeam(ctx context.Context, team model.Team) (model.Team, error) {
	sql, args, err := squirrel.
		Insert(teamTableName).
		Columns(teamColumnName).
		Values(team.Name).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.Team{}, errors.Wrap(err, "building sql")
	}

	_, err = s.getter.DefaultTrOrDB(ctx, s.pool).Exec(ctx, sql, args...)
	if constraint.IsUniqueViolation(err) {
		return model.Team{}, model.ErrTeamAlreadyExists
	}
	if err != nil {
		return model.Team{}, errors.Wrap(err, "querying sql")
	}

	return team, nil
}
