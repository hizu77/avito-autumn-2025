package main

import (
	"log"

	"github.com/hizu77/avito-autumn-2025/config"
	"github.com/hizu77/avito-autumn-2025/internal/bootstrap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize zap logger"))
	}

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("failed to initialize config", zap.Error(err))
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
		logger.Fatal("failed to initialize postgres", zap.Error(err))
	}

	app := bootstrap.InitApp(cfg, logger)

	if err := bootstrap.InitHandlers(ctx, app, pool, cfg); err != nil {
		logger.Fatal("failed to initialize handler", zap.Error(err))
	}

	if err := app.Run(ctx); err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
