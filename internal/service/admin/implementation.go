package admin

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
)

type storage interface {
	GetAdmin(ctx context.Context, id string) (model.Admin, error)
	InsertAdmin(ctx context.Context, admin model.Admin) (model.Admin, error)
}

type Service struct {
	storage   storage
	jwtSecret []byte
}

func New(
	storage storage,
	jwtSecret []byte,
) *Service {
	return &Service{
		storage:   storage,
		jwtSecret: jwtSecret,
	}
}
