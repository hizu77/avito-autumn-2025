package team

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"go.uber.org/zap"
)

type teamStorage interface {
	SaveTeam(ctx context.Context, team model.Team) (model.Team, error)
	GetTeamByName(ctx context.Context, name string) (model.Team, error)
}

type userStorage interface {
	SaveUser(ctx context.Context, user model.User) (model.User, error)
}

type transactor interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type Service struct {
	teamStorage teamStorage
	userStorage userStorage
	transactor  transactor
	logger      *zap.Logger
}

func New(
	storage teamStorage,
	userStorage userStorage,
	transactor transactor,
	logger *zap.Logger,
) *Service {
	return &Service{
		teamStorage: storage,
		userStorage: userStorage,
		transactor:  transactor,
		logger:      logger,
	}
}
