package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/swokeqq/tripngo.git/internal/domain"
	api "github.com/swokeqq/tripngo.git/internal/generated"
)

// CreateTrip — POST /api/v1/trips.
func (h *Handler) CreateTrip(w http.ResponseWriter, r *http.Request, params api.CreateTripParams) {
	var tripData api.TripData

	if err := decodeJSON(w, r, &tripData); err != nil {
		writeProblem(w, r, problemInvalidRequest, err.Error())
		return
	}
	if err := validateTripData(tripData); err != nil {
		writeProblem(w, r, problemInvalidRequest, err.Error())
		return
	}

	now := time.Now().UTC()
	trip := &domain.Trip{
		ID:       uuid.New().String(),
		UserID:   tripData.UserId.String(),
		DriverID: tripData.DriverId.String(),

		StartPoint: domain.Coordinates(tripData.StartPoint),
		EndPoint:   domain.Coordinates(tripData.EndPoint),
		Price:      tripData.Price,
		Status:     domain.StatusActive,
		StartedAt:  now,
		CreatedAt:  now,
		FinishedAt: nil,
		UpdatedAt:  &now,
	}
	history := &domain.TripStatusHistory{
		TripID:     trip.ID,
		FromStatus: nil,
		ToStatus:   string(domain.StatusActive),
		Reason:     nil,
		ChangedAt:  now,
	}

	err := h.tx.Do(r.Context(), func(ctx context.Context) error {
		if err := h.trips.Create(ctx, trip); err != nil {
			return err
		}
		return h.history.Add(ctx, history)
	})
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	tripDTO, err := toAPITrip(trip)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	w.Header().Set("Location", "/api/v1/trips/"+trip.ID)
	writeJSON(w, http.StatusCreated, tripDTO)
}

// GetTrip — GET /api/v1/trips/{tripId}.
func (h *Handler) GetTrip(w http.ResponseWriter, r *http.Request, tripId api.TripId) {
	trip, err := h.trips.GetByID(r.Context(), tripId.String())
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	tripDTO, err := toAPITrip(trip)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, tripDTO)
}

// FinishTrip — POST /api/v1/trips/{tripId}/finish.
func (h *Handler) FinishTrip(w http.ResponseWriter, r *http.Request, tripId api.TripId) {
	// TODO: реализовать
	w.WriteHeader(http.StatusNotImplemented)
}
