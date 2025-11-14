package admin

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/admin/dbmodel"
	"github.com/hizu77/avito-autumn-2025/internal/storage/common/constraint"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (s *Storage) InsertAdmin(ctx context.Context, admin model.Admin) (model.Admin, error) {
	sql, args, err := squirrel.
		Insert(adminTableName).
		Columns(columnID, columnPassword).
		Values(admin.ID, admin.PasswordHash).
		Suffix("RETURNING " + strings.Join(allColumns, ", ")).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.Admin{}, errors.Wrap(err, "building sql")
	}

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return model.Admin{}, errors.Wrap(err, "querying rows")
	}

	dbAdmin, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbmodel.Admin])
	if constraint.IsUniqueViolation(err) {
		return model.Admin{}, model.ErrAdminAlreadyExists
	}
	if err != nil {
		return model.Admin{}, errors.Wrap(err, "collecting rows")
	}

	mappedAdmin := mapDbAdminToDomainAdmin(dbAdmin)

	return mappedAdmin, nil
}
