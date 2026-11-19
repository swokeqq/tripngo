package create

import (
	"github.com/google/uuid"

	"github.com/swokeqq/tripngo.git/internal/domain/trip"
	api "github.com/swokeqq/tripngo.git/internal/generated"
	createuc "github.com/swokeqq/tripngo.git/internal/usecase/trip/create"
)

func toInput(data api.TripData) (createuc.Input, error) {
	if data.UserId == uuid.Nil {
		return createuc.Input{}, &trip.ValidationError{Field: "user_id", Reason: "can't be empty"}
	}
	if data.DriverId == uuid.Nil {
		return createuc.Input{}, &trip.ValidationError{Field: "driver_id", Reason: "can't be empty"}
	}

	input := createuc.Input{
		UserID:     data.UserId.String(),
		DriverID:   data.DriverId.String(),
		StartPoint: trip.Coordinates(data.StartPoint),
		EndPoint:   trip.Coordinates(data.EndPoint),
		Price:      data.Price,
	}

	return input, nil
}
