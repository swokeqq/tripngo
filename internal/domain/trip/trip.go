package trip

import (
	"fmt"
	"time"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
)

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Trip struct {
	ID       string
	UserID   string
	DriverID string

	StartPoint Coordinates
	EndPoint   Coordinates
	Price      int64
	Status     Status
	StartedAt  time.Time
	CreatedAt  time.Time
	FinishedAt *time.Time
	UpdatedAt  *time.Time
}

func New(
	id string,
	userID string,
	driverID string,
	startPoint Coordinates,
	endPoint Coordinates,
	price int64,
	at time.Time,
) (*Trip, error) {
	if userID == "" {
		return nil, &ValidationError{Field: "user_id", Reason: "can't be empty"}
	}
	if driverID == "" {
		return nil, &ValidationError{Field: "driver_id", Reason: "can't be empty"}
	}
	if err := validateCoordinates("start_point", startPoint); err != nil {
		return nil, err
	}

	if err := validateCoordinates("end_point", endPoint); err != nil {
		return nil, err
	}

	if price < 0 {
		return nil, &ValidationError{Field: "price", Reason: "must be non-negative"}
	}

	return &Trip{
		ID:         id,
		UserID:     userID,
		DriverID:   driverID,
		StartPoint: startPoint,
		EndPoint:   endPoint,
		Price:      price,
		Status:     StatusActive,
		StartedAt:  at,
		CreatedAt:  at,
		FinishedAt: nil,
		UpdatedAt:  &at,
	}, nil
}

func validateCoordinates(name string, coordinates Coordinates) *ValidationError {
	if coordinates.Latitude > 90 || coordinates.Latitude < -90 {
		return &ValidationError{Field: fmt.Sprintf("%s.latitude", name), Reason: "must be >= -90 and <= 90"}
	}
	if coordinates.Longitude > 180 || coordinates.Longitude < -180 {
		return &ValidationError{Field: fmt.Sprintf("%s.longitude", name), Reason: "must be >= -180 and <= 180"}
	}

	return nil
}

func (t *Trip) Finish(at time.Time) error {
	if t.Status != StatusActive {
		return ErrTripCompleted
	}
	t.Status = StatusCompleted
	t.UpdatedAt = &at
	t.FinishedAt = &at
	return nil
}
