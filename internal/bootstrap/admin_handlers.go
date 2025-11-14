package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/hizu77/avito-autumn-2025/config"
	adminhandler "github.com/hizu77/avito-autumn-2025/internal/api/admin/handler"
	admin "github.com/hizu77/avito-autumn-2025/internal/api/admin/middleware"
	adminservice "github.com/hizu77/avito-autumn-2025/internal/service/admin"
	adminstorage "github.com/hizu77/avito-autumn-2025/internal/storage/admin/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitAdminHandlers(
	app *App,
	pool *pgxpool.Pool,
	cfg *config.Config,
) error {
	secret := []byte(cfg.Secret)
	tokenAuth := jwtauth.New("HS256", secret, nil)

	storage := adminstorage.New(pool, app.logger)
	service := adminservice.New(storage, app.logger, secret)
	handler := adminhandler.New(service, app.logger)

	app.mux.Route("/admins", func(r chi.Router) {
		r.Post("/login", handler.LoginAdmin)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(admin.Authenticator)

			r.Post("/register", handler.RegisterAdmin)
		})
	})

	return nil
}
