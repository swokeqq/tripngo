package domain

import "time"

type TripStatus string

const (
	StatusActive    TripStatus = "active"
	StatusCompleted TripStatus = "completed"
)

type Coordinates struct {
	Latitude   float64
	Longtitude float64
}

type Trip struct {
	ID       string
	UserID   string
	DriverID string

	StartPoint Coordinates
	EndPoint   Coordinates
	Price      int64
	Status     TripStatus
	StartedAt  time.Time
	CreatedAt  time.Time
	FinishedAt *time.Time
	UpdatedAt  *time.Time
}
