package admin

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/common/constraint"
	"github.com/pkg/errors"
)

func (s *Storage) InsertAdmin(ctx context.Context, admin model.Admin) (model.Admin, error) {
	sql, args, err := squirrel.
		Insert(adminTableName).
		Columns(columnID, columnPassword).
		Values(admin.ID, admin.PasswordHash).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.Admin{}, errors.Wrap(err, "building sql")
	}

	_, err = s.getter.DefaultTrOrDB(ctx, s.pool).Exec(ctx, sql, args...)
	if constraint.IsUniqueViolation(err) {
		return model.Admin{}, model.ErrAdminAlreadyExists
	}
	if err != nil {
		return model.Admin{}, errors.Wrap(err, "collecting rows")
	}

	return admin, nil
}
