package team

import (
	"net/http"

	common_response "github.com/hizu77/avito-autumn-2025/internal/api/common/response"
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

func mapDomainTeamErrorToResponseErrorWithStatusCode(err error) (common_response.Error, int) {
	switch {
	case errors.Is(err, model.ErrTeamAlreadyExists):
		return common_response.NewTeamExistsError(), http.StatusBadRequest
	case errors.Is(err, model.ErrTeamDoesNotExist):
		return common_response.NewNotFoundError(), http.StatusNotFound
	default:
		return common_response.NewInternalServerError(), http.StatusInternalServerError
	}
}
