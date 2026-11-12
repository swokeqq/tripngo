package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/swokeqq/tripngo.git/internal/domain"
	api "github.com/swokeqq/tripngo.git/internal/generated"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(v)
}

func writeProblem(w http.ResponseWriter, r *http.Request, problem problemKind, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")

	w.WriteHeader(problem.status)

	instance := r.URL.Path

	prob := api.Problem{
		Code:     problem.code,
		Detail:   &detail,
		Instance: &instance,
		Status:   int32(problem.status),
		Title:    problem.title,
		Type:     problem.typ,
	}

	_ = json.NewEncoder(w).Encode(prob)
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrTripNotFound):
		writeProblem(w, r, problemTripNotFound, "Trip was not found")
	case errors.Is(err, domain.ErrTripAlreadyCompleted):
		writeProblem(w, r, problemTripCompleted, "Operation is not allowed for a completed trip")
	case errors.Is(err, domain.ErrDriverBusy):
		writeProblem(w, r, problemDriverBusy, "Driver already has an active trip")

	default:
		h.logger.ErrorContext(r.Context(), "request failed",
			slog.String("path", r.URL.Path), slog.Any("error", err))
		writeProblem(w, r, problemInternal, "Internal server error")
	}
}

// validateTripData проверяет тело createTrip по п. 5.2 ТЗ:
// user_id и driver_id не нулевые UUID, координаты в диапазонах, price >= 0.
func validateTripData(d api.TripData) error {
	// TODO: реализовать
	return errors.New("validateTripData: not implemented")
}

func toAPITrip(t *domain.Trip) (api.Trip, error) {
	driverID, err := uuid.Parse(t.DriverID)
	if err != nil {
		return api.Trip{}, fmt.Errorf("parse driverId to uuid: %w", err)
	}
	endPoint := api.Coordinates{
		Latitude:  t.EndPoint.Latitude,
		Longitude: t.EndPoint.Longitude,
	}

	id, err := uuid.Parse(t.ID)
	if err != nil {
		return api.Trip{}, fmt.Errorf("parse tripId to uuid: %w", err)
	}

	startPoint := api.Coordinates{
		Latitude:  t.StartPoint.Latitude,
		Longitude: t.StartPoint.Longitude,
	}

	userID, err := uuid.Parse(t.UserID)
	if err != nil {
		return api.Trip{}, fmt.Errorf("parse userId to uuid: %w", err)
	}

	return api.Trip{
		DriverId:       driverID,
		EndPoint:       endPoint,
		FinishedAt:     t.FinishedAt,
		Id:             id,
		LastPositionAt: nil,
		Price:          t.Price,
		StartPoint:     startPoint,
		StartedAt:      t.StartedAt,
		Status:         api.TripStatus(t.Status),
		UserId:         userID,
	}, nil
}
