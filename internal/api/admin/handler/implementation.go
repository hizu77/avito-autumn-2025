package admin

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"go.uber.org/zap"
)

type service interface {
	LoginAdmin(ctx context.Context, id string, password string) (string, error)
	RegisterAdmin(ctx context.Context, id string, password string) (model.Admin, error)
}

type Handler struct {
	service service
	logger  *zap.Logger
}

func New(
	service service,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
