package admin

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"go.uber.org/zap"
)

type storage interface {
	GetAdmin(ctx context.Context, id string) (model.Admin, error)
	InsertAdmin(ctx context.Context, admin model.Admin) (model.Admin, error)
}

type Service struct {
	storage   storage
	logger    *zap.Logger
	jwtSecret []byte
}

func New(
	storage storage,
	logger *zap.Logger,
	jwtSecret []byte,
) *Service {
	return &Service{
		storage:   storage,
		logger:    logger,
		jwtSecret: jwtSecret,
	}
}
