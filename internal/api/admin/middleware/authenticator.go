package admin

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	common_response "github.com/hizu77/avito-autumn-2025/internal/api/common/response"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil || token == nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, common_response.NewUnauthorizedError())
			return
		}

		next.ServeHTTP(w, r)
	})
}
