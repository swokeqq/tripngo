package domain

import "time"

type TripStatusHistory struct {
	ID         int64
	TripId     string
	FromStatus *string
	ToStatus   string
	Reason     *string
	ChangedAt  time.Time
}
