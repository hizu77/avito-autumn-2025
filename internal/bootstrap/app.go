package bootstrap

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hizu77/avito-autumn-2025/config"
	"github.com/hizu77/avito-autumn-2025/pkg/closer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	WriteTimeout    = 15 * time.Second
	ReadTimeout     = 15 * time.Second
	ShutdownTimeout = 15 * time.Second
)

type App struct {
	mux    *chi.Mux
	logger *zap.Logger

	httpAddr string
}

func InitApp(
	cfg *config.Config,
	logger *zap.Logger,
) *App {
	httpAddr := net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)
	mux := chi.NewRouter()

	return &App{
		mux:      mux,
		logger:   logger,
		httpAddr: httpAddr,
	}
}

func (a *App) Run(ctx context.Context) error {
	eg, _ := errgroup.WithContext(ctx)

	httpServer := &http.Server{
		Addr:         a.httpAddr,
		Handler:      a.mux,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	eg.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, "listening http")
		}

		return nil
	})

	if err := closer.AddCallback(
		CloserGroupApp,
		func() error {
			shutdownCtx, cancel := context.WithTimeout(
				context.Background(),
				ShutdownTimeout,
			)
			defer cancel()

			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				return errors.Wrap(err, "shutting down http")
			}

			a.logger.Info("http server shutdown")

			return nil
		},
	); err != nil {
		return errors.Wrap(err, "add app callback")
	}

	a.logger.Info("http server started", zap.String("http_addr", a.httpAddr))

	err := eg.Wait()
	if err != nil {
		return errors.Wrap(err, "running app")
	}

	err = closer.Wait()
	if err != nil {
		return errors.Wrap(err, "executing callbacks")
	}

	return nil
}
