package admin

import (
	"net/http"

	"github.com/hizu77/avito-autumn-2025/internal/api/admin/response"
	common_response "github.com/hizu77/avito-autumn-2025/internal/api/common/response"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
)

func mapDomainAdminErrorToResponseErrorWithStatusCode(err error) (common_response.Error, int) {
	switch {
	case errors.Is(err, model.ErrInvalidCredentials):
		return common_response.NewInvalidCredentialsError(), http.StatusUnauthorized
	case errors.Is(err, model.ErrAdminDoesNotExist):
		return common_response.NewInvalidCredentialsError(), http.StatusUnauthorized
	default:
		return common_response.NewInternalServerError(), http.StatusInternalServerError
	}
}

func mapTokenToResponseToken(token string) response.Token {
	return response.Token{
		Value: token,
	}
}

func mapTokenToResponseLoginAdmin(token string) response.LoginAdmin {
	mappedToken := mapTokenToResponseToken(token)

	return response.LoginAdmin{
		Token: mappedToken,
	}
}
