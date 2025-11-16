package bootstrap

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/pkg/closer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func InitGlobalContext(logger *zap.Logger) (context.Context, error) {
	ctx, cancel := context.WithCancel(context.Background())

	if err := closer.AddCallback(
		CloserGroupGlobalContext,
		func() error {
			logger.Info("cancel global context")
			cancel()
			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "global context callback")
	}

	return ctx, nil
}
