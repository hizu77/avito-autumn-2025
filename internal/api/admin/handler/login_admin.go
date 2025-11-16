package admin

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/admin/request"
	"github.com/hizu77/avito-autumn-2025/internal/api/httperr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h *Handler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	const op = "admin.LoginAdmin"

	var loginAdminRequest request.LoginAdmin
	if err := render.DecodeJSON(r.Body, &loginAdminRequest); err != nil {
		h.logger.Error("decoding request body",
			zap.String("op", op),
			zap.Error(err),
		)

		httperr.WriteError(w, r, httperr.CodeBadRequest, "invalid json body")
		return
	}

	if err := validateLoginAdminRequest(loginAdminRequest); err != nil {
		h.logger.Error("validating request",
			zap.String("op", op),
			zap.Error(err),
		)

		httperr.WriteError(w, r, httperr.CodeBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	token, err := h.service.LoginAdmin(
		ctx,
		loginAdminRequest.ID,
		loginAdminRequest.Password,
	)
	if err != nil {
		h.logger.Error("login admin",
			zap.String("op", op),
			zap.Error(err),
		)

		code := mapDomainAdminErrorToCode(err)
		httperr.WriteError(w, r, code)
		return
	}

	mappedToken := mapTokenToResponseLoginAdmin(token)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, mappedToken)
}

func validateLoginAdminRequest(req request.LoginAdmin) error {
	if req.ID == "" {
		return errors.New("id is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}
