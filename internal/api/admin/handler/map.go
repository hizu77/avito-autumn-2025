package admin

import (
	"net/http"

	"github.com/hizu77/avito-autumn-2025/internal/api/admin/response"
	common_response "github.com/hizu77/avito-autumn-2025/internal/api/common/response"
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

func mapDomainAdminErrorToResponseErrorWithStatusCode(err error) (common_response.Error, int) {
	switch {
	case errors.Is(err, model.ErrAdminAlreadyExists):
		return common_response.NewAdminExistsError(), http.StatusBadRequest
	case errors.Is(err, model.ErrInvalidAdminPassword):
		return common_response.NewInvalidCredentialsError(), http.StatusUnauthorized
	case errors.Is(err, model.ErrAdminDoesNotExist):
		return common_response.NewInvalidCredentialsError(), http.StatusUnauthorized
	default:
		return common_response.NewInternalServerError(), http.StatusInternalServerError
	}
}
