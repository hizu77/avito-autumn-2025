package admin

import (
	"context"
	db "database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/admin/dbmodel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (s *Storage) GetAdmin(ctx context.Context, id string) (model.Admin, error) {
	sql, args, err := squirrel.
		Select(allColumns...).
		From(adminTableName).
		Where(squirrel.Eq{columnID: id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.Admin{}, errors.Wrap(err, "building sql")
	}

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sql, args...)
	if err != nil {
		return model.Admin{}, errors.Wrap(err, "querying rows")
	}

	dbAdmin, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbmodel.Admin])
	if errors.Is(err, db.ErrNoRows) {
		return model.Admin{}, model.ErrAdminDoesNotExist
	}
	if err != nil {
		return model.Admin{}, errors.Wrap(err, "collecting rows")
	}

	mappedAdmin := mapDbAdminToDomainAdmin(dbAdmin)

	return mappedAdmin, nil
}
