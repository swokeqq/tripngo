package trip

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	tripdomain "github.com/swokeqq/tripngo.git/internal/domain/trip"
)

func (r *Repository) Finish(ctx context.Context, t *tripdomain.Trip) error {
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	query, args, err := r.qb.Update("trips").
		Set("status", tripdomain.StatusCompleted).
		Set("updated_at", t.UpdatedAt).
		Set("finished_at", t.FinishedAt).
		Where(squirrel.Eq{"id": t.ID, "status": tripdomain.StatusActive}).
		ToSql()
	if err != nil {
		return fmt.Errorf("building update query: %w", err)
	}

	db := r.db.Executor(ctx)
	tag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update trip: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return tripdomain.ErrTripCompleted
	}

	return nil
}
