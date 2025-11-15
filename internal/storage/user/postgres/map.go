package user

import (
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/user/dbmodel"
)

func mapDbUserToDomain(user dbmodel.User) model.User {
	return model.User{
		ID:       user.ID,
		Name:     user.Name,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}
