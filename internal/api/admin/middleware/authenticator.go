package admin

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/hizu77/avito-autumn-2025/internal/api/httperr"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil || token == nil {
			httperr.WriteError(w, r, httperr.CodeUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
