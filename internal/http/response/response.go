package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	tripdomain "github.com/swokeqq/tripngo.git/internal/domain/trip"
	api "github.com/swokeqq/tripngo.git/internal/generated"
	"github.com/swokeqq/tripngo.git/internal/http/domainerror"
)

const maxBodyBytes = 1 << 20

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	status, p := domainerror.Problem(err)
	path := r.URL.Path
	p.Instance = &path
	if status == http.StatusInternalServerError {
		logger.ErrorContext(r.Context(), "request failed",
			slog.String("path", r.URL.Path), slog.Any("error", err))
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(p)
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}

	return nil
}

func ToAPITrip(t *tripdomain.Trip) (api.Trip, error) {
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

	var finishedAt *time.Time
	if t.FinishedAt != nil {
		finishTime := t.FinishedAt.UTC()
		finishedAt = &finishTime
	}

	return api.Trip{
		DriverId:       driverID,
		EndPoint:       endPoint,
		FinishedAt:     finishedAt,
		Id:             id,
		LastPositionAt: nil,
		Price:          t.Price,
		StartPoint:     startPoint,
		StartedAt:      t.StartedAt.UTC(),
		Status:         api.TripStatus(t.Status),
		UserId:         userID,
	}, nil
}
