package team

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	nameQueryParam = "team_name"
)

func (h *Handler) GetTeamByName(w http.ResponseWriter, r *http.Request) {
	const op = "team.GetTeamByName"

	name := r.URL.Query().Get(nameQueryParam)

	if err := validateTeamName(name); err != nil {
		h.logger.Error("validating name",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	team, err := h.service.GetTeamByName(ctx, name)
	if err != nil {
		h.logger.Error("getting team",
			zap.String("op", op),
			zap.Error(err),
		)

		code := mapDomainTeamErrorToCode(err)
		common.WriteError(w, r, code)
		return
	}

	teamResponse := mapDomainTeamToResponseTeam(team)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, teamResponse)
}

func validateTeamName(name string) error {
	if name == "" {
		return errors.New("invalid name")
	}
	return nil
}
