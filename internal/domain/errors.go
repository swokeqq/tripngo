package domain

import "errors"

var (
	ErrTripNotFound         = errors.New("trip not found")
	ErrTripAlreadyCompleted = errors.New("trip is already completed")
)
