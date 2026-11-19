package trip

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	tripdomain "github.com/swokeqq/tripngo.git/internal/domain/trip"
)

func (r *Repository) GetByID(ctx context.Context, tripID string) (*tripdomain.Trip, error) {
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

	var trip tripdomain.Trip
	db := r.db.Executor(ctx)
	err = db.QueryRow(ctx, query, args...).Scan(
		&trip.ID, &trip.UserID, &trip.DriverID,
		&trip.StartPoint.Latitude, &trip.StartPoint.Longitude,
		&trip.EndPoint.Latitude, &trip.EndPoint.Longitude,
		&trip.Price, &trip.Status, &trip.StartedAt, &trip.CreatedAt,
		&trip.UpdatedAt, &trip.FinishedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tripdomain.ErrTripNotFound
		}
		return nil, fmt.Errorf("executing select trip query: %w", err)
	}

	return &trip, nil
}
