package team

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/common/response"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	nameQueryParam = "team_name"
)

func (i *Handler) GetTeamByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(nameQueryParam)

	if err := validateTeamName(name); err != nil {
		i.logger.Error("validating name", zap.Error(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.NewBadRequestError(err.Error()))
		return
	}

	ctx := r.Context()
	team, err := i.service.GetTeamByName(ctx, name)
	if err != nil {
		i.logger.Error("getting team", zap.Error(err))

		mappedErr, code := mapDomainTeamErrorToResponseErrorWithStatusCode(err)
		render.Status(r, code)
		render.JSON(w, r, mappedErr)
		return
	}

	teamResponse := mapDomainTeamToResponseTeam(team)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, teamResponse)
}

func validateTeamName(name string) error {
	if name == "" {
		return errors.New("invalid name")
	}

	return nil
}
