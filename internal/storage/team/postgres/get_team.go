package team

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/team/dbmodel"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (s *Storage) GetTeamByName(ctx context.Context, name string) (model.Team, error) {
	sql, args, err := squirrel.
		Expr(`
        	SELECT
            	t.name 		AS team_name,
            	u.id        AS user_id,
            	u.name      AS user_name,
            	u.is_active AS user_is_active
        	FROM teams t
        	LEFT JOIN users u ON t.name = u.team_name
        	WHERE t.name = $1`, name).
		ToSql()
	if err != nil {
		return model.Team{}, errors.Wrap(err, "building sql")
	}

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sql, args...)
	if err != nil {
		return model.Team{}, errors.Wrap(err, "fetching rows")
	}

	fetched, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Row])
	if err != nil {
		return model.Team{}, errors.Wrap(err, "collecting rows")
	}
	if len(fetched) == 0 {
		return model.Team{}, model.ErrTeamDoesNotExist
	}

	notNilUsers := collection.Filter(
		fetched,
		func(row dbmodel.Row) bool {
			return row.UID != nil
		},
	)
	mappedUsers := collection.Map(notNilUsers, mapDBRowToDomainUser)
	mappedTeam := mapDBRowToDomainTeams(fetched[0])
	mappedTeam.Members = mappedUsers

	return mappedTeam, nil
}
