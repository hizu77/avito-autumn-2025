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

func (s *Storage) GetTeamByUserID(ctx context.Context, userID string) (model.Team, error) {
	sql, args, err := squirrel.
		Expr(`
			WITH target_team AS (
		    	SELECT team_name
    			FROM users
    			WHERE id = $1
			)
			SELECT
				t.team_name  AS team_name,
    			u.id         AS user_id,
    			u.name       AS user_name,
    			u.is_active  AS user_is_active
			FROM target_team t
			JOIN users u ON u.team_name = t.team_name
			ORDER BY u.id`, userID).
		ToSql()
	if err != nil {
		return model.Team{}, errors.Wrap(err, "building sql")
	}

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sql, args...)
	if err != nil {
		return model.Team{}, errors.Wrap(err, "querying sql")
	}

	fetched, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Row])
	if err != nil {
		return model.Team{}, errors.Wrap(err, "collecting rows")
	}
	if len(fetched) == 0 {
		return model.Team{}, model.ErrTeamDoesNotExist
	}

	mappedUsers := collection.Map(fetched, mapDbRowToDomainUser)
	mappedTeam := mapDbRowToDomainTeams(fetched[0])
	mappedTeam.Members = mappedUsers

	return mappedTeam, nil
}
