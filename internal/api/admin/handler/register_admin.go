package admin

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/admin/request"
	"github.com/hizu77/avito-autumn-2025/internal/api/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h *Handler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	const op = "admin.RegisterAdmin"

	var registerAdminRequest request.RegisterAdmin
	if err := render.DecodeJSON(r.Body, &registerAdminRequest); err != nil {
		h.logger.Error("decoding request body",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, "invalid json body")
		return
	}

	if err := validateRegisterAdminRequest(registerAdminRequest); err != nil {
		h.logger.Error("validating request",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	admin, err := h.service.RegisterAdmin(
		ctx,
		registerAdminRequest.ID,
		registerAdminRequest.Password,
	)
	if err != nil {
		h.logger.Error("register admin",
			zap.String("op", op),
			zap.Error(err),
		)

		code := mapDomainAdminErrorToCode(err)
		common.WriteError(w, r, code)
		return
	}

	mappedAdmin := mapDomainAdminToResponseRegisterAdmin(admin)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, mappedAdmin)
}

func validateRegisterAdminRequest(req request.RegisterAdmin) error {
	if req.ID == "" {
		return errors.New("id is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}
