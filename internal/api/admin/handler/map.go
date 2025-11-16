package admin

import (
	"github.com/hizu77/avito-autumn-2025/internal/api/admin/response"
	"github.com/hizu77/avito-autumn-2025/internal/api/httperr"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
)

func mapTokenToResponseLoginAdmin(token string) response.LoginAdmin {
	return response.LoginAdmin{
		Token: token,
	}
}

func mapDomainAdminToResponseRegisterAdmin(admin model.Admin) response.RegisterAdmin {
	return response.RegisterAdmin{
		ID: admin.ID,
	}
}

func mapDomainAdminErrorToCode(err error) httperr.ErrorCode {
	switch {
	case errors.Is(err, model.ErrAdminAlreadyExists):
		return httperr.CodeAdminExists
	case errors.Is(err, model.ErrInvalidAdminPassword),
		errors.Is(err, model.ErrAdminDoesNotExist):
		return httperr.CodeInvalidCredentials
	default:
		return httperr.CodeInternal
	}
}
