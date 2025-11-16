package users

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"go.uber.org/zap"
)

type service interface {
	SetActive(ctx context.Context, id string, active bool) (model.User, error)
	GetUserReviewRequests(ctx context.Context, id string) ([]model.PullRequest, error)
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
