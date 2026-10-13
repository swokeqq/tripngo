package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swokeqq/tripngo.git/internal/domain"
)

type StatusHistoryRepository struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
	qb           squirrel.StatementBuilderType
}

func NewStatusHistoryRepository(db *pgxpool.Pool, queryTimeout time.Duration) *StatusHistoryRepository {
	return &StatusHistoryRepository{
		db:           db,
		queryTimeout: queryTimeout,
		qb:           squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *StatusHistoryRepository) Add(ctx context.Context, history *domain.TripStatusHistory) error {
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	query, args, err := r.qb.Insert("trip_status_history").
		Columns(
			"id", "trip_id", "from_status",
			"to_status", "reason", "changed_at",
		).
		Values(
			history.ID, history.TripId, history.FromStatus,
			history.ToStatus, history.Reason, history.ChangedAt,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("building insert query string for trip status history: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing query string for trip status history: %w", err)
	}

	return nil
}
func (r *StatusHistoryRepository) GetByTripID(ctx context.Context, tripID string) ([]domain.TripStatusHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	query, args, err := r.qb.Select(
		"id", "trip_id", "from_status",
		"to_status", "reason", "changed_at",
	).
		From("trip_status_history").
		Where(squirrel.Eq{"trip_id": tripID}).
		OrderBy("changed_at ASC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("building select query for trip status history: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for trip history status: %w", err)
	}
	defer rows.Close()

	var historyList []domain.TripStatusHistory
	for rows.Next() {
		var h domain.TripStatusHistory
		err := rows.Scan(
			&h.ID,
			&h.TripId,
			&h.FromStatus,
			&h.ToStatus,
			&h.Reason,
			&h.ChangedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning trip status history row: %w", err)
		}
		historyList = append(historyList, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterrating trip status history rows: %w", err)
	}

	return historyList, nil
}
