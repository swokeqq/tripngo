package finish

import (
	"context"
	"fmt"
	"time"

	"github.com/swokeqq/tripngo.git/internal/domain/trip"
	"github.com/swokeqq/tripngo.git/internal/model"
)

type Usecase struct {
	trips   tripRepository
	history historyRepository
	tx      txManager
	now     func() time.Time
}

func New(trips tripRepository, history historyRepository, tx txManager, now func() time.Time) *Usecase {
	return &Usecase{trips: trips, history: history, tx: tx, now: now}
}

func (uc *Usecase) Finish(ctx context.Context, tripID string) (*trip.Trip, error) {
	now := uc.now()
	var t *trip.Trip
	statusActive := string(trip.StatusActive)

	err := uc.tx.Do(ctx, func(ctx context.Context) error {
		found, getErr := uc.trips.GetByID(ctx, tripID)
		if getErr != nil {
			return getErr
		}
		t = found

		if finErr := t.Finish(now); finErr != nil {
			return finErr
		}
		if saveErr := uc.trips.Finish(ctx, t); saveErr != nil {
			return fmt.Errorf("save trip: %w", saveErr)
		}
		if comErr := uc.history.Add(ctx, &model.TripStatusHistory{
			TripID:     tripID,
			FromStatus: &statusActive,
			ToStatus:   string(trip.StatusCompleted),
			ChangedAt:  now,
		}); comErr != nil {
			return fmt.Errorf("add to history: %w", comErr)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}
