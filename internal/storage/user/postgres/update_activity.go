package user

import (
	"context"
	db "database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/user/dbmodel"
	"github.com/pkg/errors"
)

func (s *Storage) UpdateActivity(ctx context.Context, id string, activity bool) (model.User, error) {
	sql, args, err := squirrel.
		Update(tableName).
		Set(columnIsActive, activity).
		Where(squirrel.Eq{columnID: id}).
		Suffix("RETURNING " + strings.Join(allColumns, ", ")).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.User{}, errors.Wrap(err, "building sql")
	}

	var dbUser dbmodel.User
	err = s.getter.DefaultTrOrDB(ctx, s.pool).
		QueryRow(ctx, sql, args...).
		Scan(
			&dbUser.ID,
			&dbUser.Name,
			&dbUser.TeamName,
			&dbUser.IsActive,
		)
	if errors.Is(err, db.ErrNoRows) {
		return model.User{}, model.ErrUserDoesNotExist
	}
	if err != nil {
		return model.User{}, errors.Wrap(err, "fetching row")
	}

	mappedUser := mapDBUserToDomain(dbUser)

	return mappedUser, nil
}
