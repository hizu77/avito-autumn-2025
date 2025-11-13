package team

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/common/response"
	"github.com/hizu77/avito-autumn-2025/internal/api/team/request"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (i *Handler) SaveTeam(w http.ResponseWriter, r *http.Request) {
	var saveTeamRequest request.SaveTeam
	if err := render.DecodeJSON(r.Body, &saveTeamRequest); err != nil {
		i.logger.Error("decoding request body", zap.Error(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.NewBadRequestError("invalid json body"))
		return
	}

	if err := validateSaveTeamRequest(saveTeamRequest); err != nil {
		i.logger.Error("validating request", zap.Error(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.NewBadRequestError(err.Error()))
		return
	}

	ctx := r.Context()
	mappedTeam := mapRequestSaveTeamToDomainTeam(saveTeamRequest)
	team, err := i.service.SaveTeam(ctx, mappedTeam)
	if err != nil {
		i.logger.Error("saving team", zap.Error(err))

		mappedErr, code := mapDomainTeamErrorToResponseErrorWithStatusCode(err)
		render.Status(r, code)
		render.JSON(w, r, mappedErr)
		return
	}

	teamResponse := mapDomainTeamToResponseSaveTeam(team)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, teamResponse)
}

func validateSaveTeamRequest(saveTeamRequest request.SaveTeam) error {
	if saveTeamRequest.Name == "" {
		return errors.New("name is required")
	}

	if len(saveTeamRequest.Members) == 0 {
		return errors.New("members are required")
	}

	return nil
}
