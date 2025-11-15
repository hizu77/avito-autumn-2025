package admin

import (
	"github.com/hizu77/avito-autumn-2025/internal/api/admin/response"
	"github.com/hizu77/avito-autumn-2025/internal/api/common"
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

func mapDomainAdminErrorToCode(err error) common.ErrorCode {
	switch {
	case errors.Is(err, model.ErrAdminAlreadyExists):
		return common.CodeAdminExists
	case errors.Is(err, model.ErrInvalidAdminPassword),
		errors.Is(err, model.ErrAdminDoesNotExist):
		return common.CodeInvalidCredentials
	default:
		return common.CodeInternal
	}
}
