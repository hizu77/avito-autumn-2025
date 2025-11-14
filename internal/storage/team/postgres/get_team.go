package team

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/team/dbmodel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (s *Storage) GetTeam(ctx context.Context, id string) (model.Team, error) {
	sql, args, err := squirrel.
		Select(
			"t."+teamColumnName,
			"u."+userColumnID,
			"u."+userColumnName,
			"u."+userColumnIsActive,
		).
		From(teamTableName + " t").
		Join(userColumnTeamName + " u ON " + "t." + teamColumnName + " = u." + userColumnTeamName).
		Where(squirrel.Eq{"t." + teamColumnName: id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.Team{}, errors.Wrap(err, "building sql")
	}

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sql, args...)
	if err != nil {
		return model.Team{}, errors.Wrap(err, "fetching rows")
	}

	fetched, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Row])
	if len(fetched) == 0 {
		return model.Team{}, model.ErrTeamDoesNotExist
	}
	if err != nil {
		return model.Team{}, errors.Wrap(err, "collecting rows")
	}

	mappedUsers := mapDbRowsToDomainUsers(fetched)
	mappedTeam := mapDbRowToDomainTeams(fetched[0])
	mappedTeam.Members = mappedUsers

	return mappedTeam, nil
}
