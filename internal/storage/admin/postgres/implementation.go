package admin

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func New(
	pool *pgxpool.Pool,
	logger *zap.Logger,
) *Storage {
	return &Storage{
		pool:   pool,
		logger: logger,
	}
}
