package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/avito-autumn-2025/internal/model"
)

func (s *Storage) UpdateActivity(ctx context.Context, id string, activity bool) (model.User, error) {
	sql, args, err := squirrel.
		
}