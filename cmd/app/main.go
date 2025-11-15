package main

import (
	"log"

	"github.com/hizu77/avito-autumn-2025/config"
	"github.com/hizu77/avito-autumn-2025/db"
	"github.com/hizu77/avito-autumn-2025/internal/bootstrap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// TODO сделать норм логи
// TODO вынести накат миграцй и создание дефолтного админа в сайдкар
// TODO чекнуть как можно лучше сделать ErrorResponse
// TODO чекнуть RETURNING, где он не нужен убрать его
// TODO транзакции — ответственность сервисов, а не репозиториев

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize zap logger"))
	}

	cfg, err := config.New(logger)
	if err != nil {
		log.Fatal(errors.Wrap(err, "init config"))
	}

	bootstrap.InitCloser()

	ctx, err := bootstrap.InitGlobalContext(logger)
	if err != nil {
		logger.Fatal("failed to initialize global context", zap.Error(err))
	}

	pool, err := bootstrap.InitPostgres(
		ctx,
		cfg.Postgres.URL,
		logger,
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "init postgres"))
	}

	err = db.Migrate(pool)
	if err != nil {
		log.Fatal(errors.Wrap(err, "migrating"))
	}

	app := bootstrap.InitApp(cfg, logger)

	if err := bootstrap.InitHandlers(app, pool, cfg); err != nil {
		log.Fatal(errors.Wrap(err, "init admin handlers"))
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal(errors.Wrap(err, "running app"))
	}
}
