package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	api "github.com/swokeqq/tripngo.git/internal/generated"
	"github.com/swokeqq/tripngo.git/internal/repository/postgres"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Handler struct {
	api.Unimplemented
	dbPool       *pgxpool.Pool
	queryTimeout time.Duration
	logger       *slog.Logger

	trips   *postgres.TripRepository
	history *postgres.StatusHistoryRepository
	tx      TxManager
}

func New(
	dbPool *pgxpool.Pool,
	queryTimeout time.Duration,
	logger *slog.Logger,
	trips *postgres.TripRepository,
	history *postgres.StatusHistoryRepository,
	tx TxManager,
) *Handler {
	return &Handler{
		dbPool:       dbPool,
		queryTimeout: queryTimeout,
		logger:       logger,
		trips:        trips,
		history:      history,
		tx:           tx,
	}
}
