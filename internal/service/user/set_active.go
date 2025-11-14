package user

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
)

func (s *Service) SetActive(ctx context.Context, id string, active bool) (model.User, error) {
	return s.userStorage.UpdateActivity(ctx, id, active)
}
