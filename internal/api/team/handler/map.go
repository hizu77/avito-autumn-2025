package team

import (
	"github.com/hizu77/avito-autumn-2025/internal/api/httperr"
	"github.com/hizu77/avito-autumn-2025/internal/api/team/request"
	"github.com/hizu77/avito-autumn-2025/internal/api/team/response"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/pkg/errors"
)

func mapRequestMemberToDomainUser(member request.TeamMember) model.User {
	return model.User{
		ID:       member.ID,
		Name:     member.Name,
		IsActive: member.IsActive,
	}
}

func mapRequestSaveTeamToDomainTeam(team request.SaveTeam) model.Team {
	mappedUsers := collection.Map(team.Members, mapRequestMemberToDomainUser)

	return model.Team{
		Name:    team.Name,
		Members: mappedUsers,
	}
}

func mapDomainUserToResponseMember(member model.User) response.TeamMember {
	return response.TeamMember{
		ID:       member.ID,
		Name:     member.Name,
		IsActive: member.IsActive,
	}
}

func mapDomainTeamToResponseTeam(team model.Team) response.Team {
	mappedUsers := collection.Map(team.Members, mapDomainUserToResponseMember)

	return response.Team{
		Name:    team.Name,
		Members: mappedUsers,
	}
}

func mapDomainTeamToResponseSaveTeam(team model.Team) response.SaveTeam {
	mappedTeam := mapDomainTeamToResponseTeam(team)

	return response.SaveTeam{
		Team: mappedTeam,
	}
}

func mapDomainTeamErrorToCode(err error) httperr.ErrorCode {
	switch {
	case errors.Is(err, model.ErrTeamAlreadyExists):
		return httperr.CodeTeamExists
	case errors.Is(err, model.ErrTeamDoesNotExist):
		return httperr.CodeNotFound
	default:
		return httperr.CodeInternal
	}
}
