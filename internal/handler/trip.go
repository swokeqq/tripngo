package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	domain "github.com/swokeqq/tripngo.git/internal/domain"
	api "github.com/swokeqq/tripngo.git/internal/generated"
	"github.com/swokeqq/tripngo.git/internal/model"
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
	history := &model.TripStatusHistory{
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
	statusActive := string(domain.StatusActive)
	now := time.Now().UTC()
	tripID := tripId.String()
	var trip *domain.Trip

	history := &model.TripStatusHistory{
		TripID:     tripID,
		FromStatus: &statusActive,
		ToStatus:   string(domain.StatusCompleted),
		ChangedAt:  now,
	}

	err := h.tx.Do(r.Context(), func(ctx context.Context) error {
		var err error
		trip, err = h.trips.Finish(ctx, tripID, now)
		if err != nil {
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

	writeJSON(w, http.StatusOK, tripDTO)
}

func ParamErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var invalidParam *api.InvalidParamFormatError
	if errors.As(err, &invalidParam) {
		detail := fmt.Sprintf("%s has invalid format", invalidParam.ParamName)
		writeProblem(w, r, problemInvalidRequest, detail)
		return
	}

	writeProblem(w, r, problemInvalidRequest, "invalid request parameters")
}
