package create

import (
	"context"
	"fmt"
	"time"

	"github.com/swokeqq/tripngo.git/internal/domain/trip"
	"github.com/swokeqq/tripngo.git/internal/model"
)

type Input struct {
	UserID     string
	DriverID   string
	StartPoint trip.Coordinates
	EndPoint   trip.Coordinates
	Price      int64
}

type Usecase struct {
	trips   tripRepository
	history historyRepository
	tx      txManager
	now     func() time.Time
	newID   func() string
}

func New(trips tripRepository, history historyRepository,
	tx txManager, now func() time.Time, newID func() string,
) *Usecase {
	return &Usecase{trips: trips, history: history, tx: tx, now: now, newID: newID}
}

func (uc *Usecase) Create(ctx context.Context, in Input) (*trip.Trip, error) {
	now := uc.now()
	id := uc.newID()

	t, err := trip.New(id, in.UserID, in.DriverID, in.StartPoint, in.EndPoint, in.Price, now)
	if err != nil {
		return nil, err
	}

	if err = uc.tx.Do(ctx, func(ctx context.Context) error {
		if saveErr := uc.trips.Save(ctx, t); saveErr != nil {
			return fmt.Errorf("save trip: %w", saveErr)
		}
		if addErr := uc.history.Add(ctx, &model.TripStatusHistory{
			TripID: t.ID, FromStatus: nil, ToStatus: string(trip.StatusActive), ChangedAt: now,
		}); addErr != nil {
			return fmt.Errorf("write history: %w", addErr)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return t, nil
}
