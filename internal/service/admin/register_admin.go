package admin

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterAdmin(
	ctx context.Context,
	id string,
	password string,
) (model.Admin, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.Admin{}, errors.Wrap(err, "generating hash")
	}

	admin := model.Admin{
		ID:           id,
		PasswordHash: string(hash),
	}

	return s.storage.InsertAdmin(ctx, admin)
}
