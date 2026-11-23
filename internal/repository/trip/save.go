package trip

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	tripdomain "github.com/swokeqq/tripngo.git/internal/domain/trip"
)

func (r *Repository) Save(ctx context.Context, trip *tripdomain.Trip) error {
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	query, args, err := r.qb.Insert("trips").
		Columns(
			"id", "user_id", "driver_id",
			"start_latitude", "start_longitude",
			"end_latitude", "end_longitude",
			"price", "status", "started_at",
			"created_at", "updated_at",
		).
		Values(
			trip.ID, trip.UserID, trip.DriverID,
			trip.StartPoint.Latitude, trip.StartPoint.Longitude,
			trip.EndPoint.Latitude, trip.EndPoint.Longitude,
			trip.Price, trip.Status, trip.StartedAt,
			trip.CreatedAt, trip.UpdatedAt,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for trip: %w", err)
	}

	db := r.db.Executor(ctx)
	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == pgUniqueViolation &&
			pgErr.ConstraintName == "trips_driver_active_unique_idx" {
			return tripdomain.ErrDriverBusy
		}
		return fmt.Errorf("executing insert trip query: %w", err)
	}

	return nil
}
