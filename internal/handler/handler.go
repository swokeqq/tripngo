package handler

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	api "github.com/swokeqq/tripngo.git/internal/generated"
)

type Handler struct {
	api.Unimplemented
	dbPool       *pgxpool.Pool
	queryTimeout time.Duration
}

func New(dbPool *pgxpool.Pool, queryTimeout time.Duration) *Handler {
	return &Handler{
		dbPool:       dbPool,
		queryTimeout: queryTimeout,
	}
}
