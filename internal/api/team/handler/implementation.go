package team

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"go.uber.org/zap"
)

type service interface {
	SaveTeam(ctx context.Context, team model.Team) (model.Team, error)
	GetTeamByName(ctx context.Context, name string) (model.Team, error)
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
