package users

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/common/response"
	"github.com/hizu77/avito-autumn-2025/internal/api/user/request"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h *Handler) SetActive(w http.ResponseWriter, r *http.Request) {
	var setActiveRequest request.SetActive
	if err := render.DecodeJSON(r.Body, &setActiveRequest); err != nil {
		h.logger.Error("decoding request body", zap.Error(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.NewBadRequestError("invalid json body"))
		return
	}

	if err := validateSetActiveRequest(setActiveRequest); err != nil {
		h.logger.Error("validating request", zap.Error(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.NewBadRequestError(err.Error()))
		return
	}

	ctx := r.Context()
	user, err := h.service.SetActive(ctx, setActiveRequest.ID, setActiveRequest.IsActive)
	if err != nil {
		h.logger.Error("setting active user", zap.Error(err))

		mappedErr, code := mapDomainUserErrorToResponseErrorWithStatusCode(err)
		render.Status(r, code)
		render.JSON(w, r, mappedErr)
		return
	}

	userResponse := mapDomainUserToResponseSetActive(user)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, userResponse)
}

func validateSetActiveRequest(setActiveRequest request.SetActive) error {
	if setActiveRequest.ID == "" {
		return errors.New("id is required")
	}

	return nil
}
