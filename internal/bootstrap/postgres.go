package bootstrap

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/pkg/closer"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func InitPostgres(
	ctx context.Context,
	connectionString string,
	logger *zap.Logger,
) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(
		ctx,
		connectionString,
	)
	if err != nil {
		return nil, errors.Wrap(err, "init postgres")
	}

	if err := closer.AddCallback(
		CloserGroupConnections,
		func() error {
			logger.Info("closing database connection")
			pool.Close()
			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "postgres callback")
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "readiness probe")
	}

	return pool, nil
}
