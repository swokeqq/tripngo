package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swokeqq/tripngo.git/internal/domain"
)

type TripRepository struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
	qb           squirrel.StatementBuilderType
}

func NewTripRepository(db *pgxpool.Pool, queryTimeout time.Duration) *TripRepository {
	return &TripRepository{
		db:           db,
		queryTimeout: queryTimeout,
		qb:           squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *TripRepository) Create(ctx context.Context, trip *domain.Trip) error {
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	query, args, err := r.qb.Insert("trips").
		Columns(
			"id", "user_id", "driver_id",
			"start_latitude", "start_longitude",
			"end_latitude", "end_longitude",
			"price", "status", "started_at", "created_at",
		).
		Values(
			trip.ID, trip.UserID, trip.DriverID,
			trip.StartPoint.Latitude, trip.StartPoint.Longitude,
			trip.EndPoint.Latitude, trip.EndPoint.Longitude,
			trip.Price, trip.Status, trip.StartedAt, trip.CreatedAt,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("building insert query for trip: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing insert trip query: %w", err)
	}

	return nil
}

func (r *TripRepository) GetByID(ctx context.Context, tripID string) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	query, args, err := r.qb.Select(
		"id", "user_id", "driver_id",
		"start_latitude", "start_longitude",
		"end_latitude", "end_longitude",
		"price", "status", "started_at", "created_at",
		"updated_at", "finished_at",
	).
		From("trips").
		Where(squirrel.Eq{"id": tripID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("building select query for trip: %w", err)
	}

	var trip domain.Trip
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&trip.ID, &trip.UserID, &trip.DriverID,
		&trip.StartPoint.Latitude, &trip.StartPoint.Longitude,
		&trip.EndPoint.Latitude, &trip.EndPoint.Longitude,
		&trip.Price, &trip.Status, &trip.StartedAt, &trip.CreatedAt,
		&trip.UpdatedAt, &trip.FinishedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTripNotFound
		}
		return nil, fmt.Errorf("executing select trip query: %w", err)
	}

	return &trip, nil
}
