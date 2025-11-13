package main

import (
	"log"

	"github.com/hizu77/avito-autumn-2025/internal/bootstrap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize zap logger"))
	}

	bootstrap.InitCloser()

	ctx, err := bootstrap.InitGlobalContext(logger)
	if err != nil {
		logger.Fatal("failed to initialize global context", zap.Error(err))
	}
}
