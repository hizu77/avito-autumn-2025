package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/pkg/errors"
)

func (s *Storage) SaveUsers(ctx context.Context, users []model.User) ([]model.User, error) {
	if len(users) == 0 {
		return nil, nil
	}

	ids := collection.Map(users, model.User.GetID)
	names := collection.Map(users, model.User.GetName)
	teamNames := collection.Map(users, model.User.GetTeamName)
	actives := collection.Map(users, model.User.GetIsActive)

	sql, args, err := squirrel.
		Expr(`
            INSERT INTO users (id, name, team_name, is_active)
            SELECT * FROM unnest(
                $1::text[],
                $2::text[],
                $3::text[],
                $4::bool[]
            ) AS t(id, name, team_name, is_active)
            ON CONFLICT (id) DO UPDATE
            SET team_name = EXCLUDED.team_name`,
			ids, names, teamNames, actives).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	_, err = s.getter.DefaultTrOrDB(ctx, s.pool).Exec(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing sql")
	}

	return users, nil
}
