package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swokeqq/tripngo.git/internal/config"
	"github.com/swokeqq/tripngo.git/internal/handler"
	"github.com/swokeqq/tripngo.git/internal/repository/postgres"

	api "github.com/swokeqq/tripngo.git/internal/generated"
)

type App struct {
	httpServer *http.Server
	dbPool     *pgxpool.Pool
	logger     *slog.Logger
}

func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	dbPool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("initializing postgresql pool: %w", err)
	}

	h := handler.New(dbPool, cfg.Database.QueryTimeout)

	r := chi.NewRouter()

	api.HandlerFromMux(h, r)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadTimeout:       cfg.HTTPReadTimeout,
		ReadHeaderTimeout: cfg.HTTPReadHeaderTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       cfg.HTTPIdleTimeout,
	}

	return &App{
		httpServer: srv,
		dbPool:     dbPool,
		logger:     logger,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info("starting HTTP server", slog.String("addr:", a.httpServer.Addr))
	if err := a.httpServer.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("stopping HTTP server...")
	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Error("http server shutdown failed", slog.Any("error", err))
	}

	a.dbPool.Close()
	a.logger.Info("postgresql pool closed successfully")

	return nil
}
