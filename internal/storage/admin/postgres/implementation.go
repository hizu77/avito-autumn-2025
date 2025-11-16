package admin

import (
	"github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool   *pgxpool.Pool
	getter *pgxv5.CtxGetter
}

func New(
	pool *pgxpool.Pool,
	getter *pgxv5.CtxGetter,
) *Storage {
	return &Storage{
		pool:   pool,
		getter: getter,
	}
}
