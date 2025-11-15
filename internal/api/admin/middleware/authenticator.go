package admin

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/hizu77/avito-autumn-2025/internal/api/common"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil || token == nil {
			common.WriteError(w, r, common.CodeUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
