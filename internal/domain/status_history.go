package domain

import "time"

type TripStatusHistory struct {
	ID         string
	TripId     string
	FromStatus *string
	ToStatus   string
	Reason     *string
	ChangedAt  time.Time
}
