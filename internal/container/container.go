package container

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/swokeqq/tripngo.git/internal/config"
	api "github.com/swokeqq/tripngo.git/internal/generated"
	"github.com/swokeqq/tripngo.git/internal/http/ready"
	"github.com/swokeqq/tripngo.git/internal/http/server"
	createh "github.com/swokeqq/tripngo.git/internal/http/trip/create"
	finishh "github.com/swokeqq/tripngo.git/internal/http/trip/finish"
	geth "github.com/swokeqq/tripngo.git/internal/http/trip/get"
	"github.com/swokeqq/tripngo.git/internal/infra/postgres"
	"github.com/swokeqq/tripngo.git/internal/infra/postgres/transaction"
	triprepo "github.com/swokeqq/tripngo.git/internal/repository/trip"
	createuc "github.com/swokeqq/tripngo.git/internal/usecase/trip/create"
	finishuc "github.com/swokeqq/tripngo.git/internal/usecase/trip/finish"
	getuc "github.com/swokeqq/tripngo.git/internal/usecase/trip/get"
)

type Container struct {
	httpServer *http.Server
	dbPool     *pgxpool.Pool
	logger     *slog.Logger
}

func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Container, error) {
	dbPool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("initializing postgresql pool: %w", err)
	}

	txManager := transaction.NewTxManager(dbPool)
	trips := triprepo.NewRepository(txManager, cfg.Database.QueryTimeout)
	history := triprepo.NewHistoryRepository(txManager, cfg.Database.QueryTimeout)

	now := func() time.Time { return time.Now().UTC() }

	createUC := createuc.New(trips, history, txManager, now, uuid.NewString)
	getUC := getuc.New(trips)
	finishUC := finishuc.New(trips, history, txManager, now)

	createHandler := createh.New(createUC, logger)
	getHandler := geth.New(getUC, logger)
	finishHandler := finishh.New(finishUC, logger)
	readyHandler := ready.New(dbPool, cfg.Database.QueryTimeout)

	apiSrv := server.New(createHandler, getHandler, finishHandler, readyHandler, logger)

	r := chi.NewRouter()

	api.HandlerWithOptions(apiSrv, api.ChiServerOptions{
		BaseRouter:       r,
		ErrorHandlerFunc: apiSrv.ParamErrorHandler,
	})

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadTimeout:       cfg.HTTPReadTimeout,
		ReadHeaderTimeout: cfg.HTTPReadHeaderTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       cfg.HTTPIdleTimeout,
	}

	return &Container{
		httpServer: srv,
		dbPool:     dbPool,
		logger:     logger,
	}, nil
}

func (c *Container) Run() error {
	c.logger.Info("starting HTTP server", slog.String("addr", c.httpServer.Addr))
	if err := c.httpServer.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

func (c *Container) Stop(ctx context.Context) error {
	c.logger.Info("stopping HTTP server...")
	err := c.httpServer.Shutdown(ctx)
	c.dbPool.Close()
	if err != nil {
		return fmt.Errorf("http server shutdown failed: %w", err)
	}

	c.logger.Info("postgresql pool closed successfully")
	return nil
}
