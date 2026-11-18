package get

import (
	"context"

	"github.com/swokeqq/tripngo.git/internal/domain/trip"
)

type Usecase struct {
	trips tripRepository
}

func New(trips tripRepository) *Usecase {
	return &Usecase{trips: trips}
}

func (uc *Usecase) GetByID(ctx context.Context, id string) (*trip.Trip, error) {
	return uc.trips.GetByID(ctx, id)
}
