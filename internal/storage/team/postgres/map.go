package team

import (
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/team/dbmodel"
)

func mapDbRowToDomainUser(row dbmodel.Row) model.User {
	return model.User{
		ID:       row.UID,
		Name:     row.UName,
		TeamName: row.TName,
		IsActive: row.UIsActive,
	}
}

func mapDbRowToDomainTeams(row dbmodel.Row) model.Team {
	return model.Team{
		Name: row.TName,
	}
}
