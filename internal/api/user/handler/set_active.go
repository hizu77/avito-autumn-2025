package users

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/common"
	"github.com/hizu77/avito-autumn-2025/internal/api/user/request"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h *Handler) SetActive(w http.ResponseWriter, r *http.Request) {
	const op = "users.SetActive"

	var setActiveRequest request.SetActive
	if err := render.DecodeJSON(r.Body, &setActiveRequest); err != nil {
		h.logger.Error("decoding request body",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, "invalid json body")
		return
	}

	if err := validateSetActiveRequest(setActiveRequest); err != nil {
		h.logger.Error("validating request",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	user, err := h.service.SetActive(ctx, setActiveRequest.ID, setActiveRequest.IsActive)
	if err != nil {
		h.logger.Error("setting active user",
			zap.String("op", op),
			zap.Error(err),
		)

		code := mapDomainUserErrorToCode(err)
		common.WriteError(w, r, code)
		return
	}

	userResponse := mapDomainUserToResponseSetActive(user)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, userResponse)
}

func validateSetActiveRequest(req request.SetActive) error {
	if req.ID == "" {
		return errors.New("id is required")
	}
	return nil
}
