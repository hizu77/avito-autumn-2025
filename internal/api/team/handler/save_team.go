package team

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/common"
	"github.com/hizu77/avito-autumn-2025/internal/api/team/request"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h *Handler) SaveTeam(w http.ResponseWriter, r *http.Request) {
	const op = "team.SaveTeam"

	var saveTeamRequest request.SaveTeam
	if err := render.DecodeJSON(r.Body, &saveTeamRequest); err != nil {
		h.logger.Error("decoding request body",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, "invalid json body")
		return
	}

	if err := validateSaveTeamRequest(saveTeamRequest); err != nil {
		h.logger.Error("validating request",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	mappedTeam := mapRequestSaveTeamToDomainTeam(saveTeamRequest)

	team, err := h.service.SaveTeam(ctx, mappedTeam)
	if err != nil {
		h.logger.Error("saving team",
			zap.String("op", op),
			zap.Error(err),
		)

		code := mapDomainTeamErrorToCode(err)
		common.WriteError(w, r, code)
		return
	}

	teamResponse := mapDomainTeamToResponseSaveTeam(team)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, teamResponse)
}

func validateSaveTeamRequest(req request.SaveTeam) error {
	if req.Name == "" {
		return errors.New("name is required")
	}

	if len(req.Members) == 0 {
		return errors.New("members are required")
	}

	return nil
}
