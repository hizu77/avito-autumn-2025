package team

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
)

func (s *Service) GetTeamByName(ctx context.Context, name string) (model.Team, error) {
	return s.teamStorage.GetTeamByName(ctx, name)
}
