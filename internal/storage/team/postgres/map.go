package team

import (
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/team/dbmodel"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
)

func mapDbRowsToDomainUsers(rows []dbmodel.Row) []model.User {
	return collection.Reduce(
		rows,
		func(users []model.User, row dbmodel.Row) []model.User {
			user := model.User{
				ID:       row.UID,
				Name:     row.UName,
				TeamName: row.UTeamName,
				IsActive: row.UIsActive,
			}

			users = append(users, user)

			return users
		},
		[]model.User{},
	)
}

func mapDbRowToDomainTeams(row dbmodel.Row) model.Team {
	return model.Team{
		Name: row.TName,
	}
}
