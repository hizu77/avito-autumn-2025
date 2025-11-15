package bootstrap

import (
	"github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/hizu77/avito-autumn-2025/config"
	adminhandler "github.com/hizu77/avito-autumn-2025/internal/api/admin/handler"
	middleware "github.com/hizu77/avito-autumn-2025/internal/api/admin/middleware"
	pullrequesthandler "github.com/hizu77/avito-autumn-2025/internal/api/pull_requests/handler"
	teamhandler "github.com/hizu77/avito-autumn-2025/internal/api/team/handler"
	userhandler "github.com/hizu77/avito-autumn-2025/internal/api/user/handler"
	adminservice "github.com/hizu77/avito-autumn-2025/internal/service/admin"
	pullrequestservice "github.com/hizu77/avito-autumn-2025/internal/service/pull_request"
	teamservice "github.com/hizu77/avito-autumn-2025/internal/service/team"
	userservice "github.com/hizu77/avito-autumn-2025/internal/service/user"
	adminstorage "github.com/hizu77/avito-autumn-2025/internal/storage/admin/postgres"
	pullrequeststorage "github.com/hizu77/avito-autumn-2025/internal/storage/pull_request/postgres"
	teamstorage "github.com/hizu77/avito-autumn-2025/internal/storage/team/postgres"
	userstorage "github.com/hizu77/avito-autumn-2025/internal/storage/user/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitHandlers(
	app *App,
	pool *pgxpool.Pool,
	cfg *config.Config,
) {
	secret := []byte(cfg.Secret)
	tokenAuth := jwtauth.New("HS256", secret, nil)

	trManager := manager.Must(pgxv5.NewDefaultFactory(pool))
	trGetter := pgxv5.DefaultCtxGetter

	adminStorage := adminstorage.New(pool, trGetter)
	userStorage := userstorage.New(pool, trGetter)
	teamStorage := teamstorage.New(pool, trGetter)
	pullRequestStorage := pullrequeststorage.New(pool, trGetter)

	adminService := adminservice.New(adminStorage, secret)
	userService := userservice.New(userStorage, pullRequestStorage)
	teamService := teamservice.New(userStorage, teamStorage, trManager)
	pullRequestService := pullrequestservice.New(
		teamStorage,
		pullRequestStorage,
		trManager,
	)

	adminHandler := adminhandler.New(adminService, app.logger)
	userHandler := userhandler.New(userService, app.logger)
	teamHandler := teamhandler.New(teamService, app.logger)
	pullRequestHandler := pullrequesthandler.New(pullRequestService, app.logger)

	app.mux.Route("/admins", func(r chi.Router) {
		r.Post("/login", adminHandler.LoginAdmin)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(middleware.Authenticator)
			r.Post("/register", adminHandler.RegisterAdmin)
		})
	})

	app.mux.Route("/team", func(r chi.Router) {
		r.Post("/add", teamHandler.SaveTeam)
		r.Get("/get", teamHandler.GetTeamByName)
	})

	app.mux.Route("/users", func(r chi.Router) {
		r.Get("/getReview", userHandler.GetUserReviewRequests)
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(middleware.Authenticator)
			r.Post("/setIsActive", userHandler.SetActive)
		})
	})

	app.mux.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", pullRequestHandler.CreatePullRequest)
		r.Post("/merge", pullRequestHandler.MergePullRequest)
		r.Post("/reassign", pullRequestHandler.ReassignPullRequest)
	})
}
