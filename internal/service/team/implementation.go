package team

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/hizu77/avito-autumn-2025/internal/model"
)

//go:generate mockgen -source=implementation.go -destination=../../mock/team/storage.go -package=mock -mock_names teamStorage=TeamStorage,userStorage=UserStorage
type (
	userStorage interface {
		SaveUsers(ctx context.Context, users []model.User) ([]model.User, error)
	}

	teamStorage interface {
		SaveTeam(ctx context.Context, team model.Team) (model.Team, error)
		GetTeamByName(ctx context.Context, name string) (model.Team, error)
	}
)

type Service struct {
	userStorage userStorage
	teamStorage teamStorage

	trManager trm.Manager
}

func New(
	userStorage userStorage,
	teamStorage teamStorage,
	trManager trm.Manager,
) *Service {
	return &Service{
		userStorage: userStorage,
		teamStorage: teamStorage,
		trManager:   trManager,
	}
}
