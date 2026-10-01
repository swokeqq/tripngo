package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/swokeqq/tripngo.git/internal/config"
	"github.com/swokeqq/tripngo.git/internal/handler"
)

type App struct {
	httpServer *http.Server
}

func New(cfg *config.Config, h *handler.Handler) *App {
	r := chi.NewRouter()
	srv := http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadTimeout:       cfg.HTTPReadTimeout,
		ReadHeaderTimeout: cfg.HTTPReadHeaderTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       cfg.HTTPIdleTimeout,
	}

	return &App{
		httpServer: &srv,
	}
}

func (a *App) Run() error {
	fmt.Printf("starting HTTP server on %s", a.httpServer.Addr)
	if err := a.httpServer.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and server: %w", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	fmt.Printf("stopping HTTP server...")
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}
	return nil
}
