package pullrequest

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"go.uber.org/zap"
)

type service interface {
	CreatePullRequest(ctx context.Context, request model.PullRequest) (model.PullRequest, error)
	MergePullRequest(ctx context.Context, id string) (model.PullRequest, error)
	ReassignPullRequest(ctx context.Context, id string, reviewerID string) (model.PullRequest, error)
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
