package create

import (
	"log/slog"
	"net/http"

	"github.com/swokeqq/tripngo.git/internal/domain/trip"
	api "github.com/swokeqq/tripngo.git/internal/generated"
	"github.com/swokeqq/tripngo.git/internal/http/response"
)

type Handler struct {
	uc     usecase
	logger *slog.Logger
}

func New(uc usecase, logger *slog.Logger) *Handler {
	return &Handler{uc: uc, logger: logger}
}

func (h *Handler) CreateTrip(w http.ResponseWriter, r *http.Request, _ api.CreateTripParams) {
	var data api.TripData
	if err := response.DecodeJSON(w, r, &data); err != nil {
		response.WriteError(w, r, h.logger, &trip.ValidationError{Field: "body", Reason: err.Error()})
		return
	}

	input, err := toInput(data)
	if err != nil {
		response.WriteError(w, r, h.logger, err)
		return
	}

	t, err := h.uc.Create(r.Context(), input)
	if err != nil {
		response.WriteError(w, r, h.logger, err)
		return
	}

	tripDTO, err := response.ToAPITrip(t)
	if err != nil {
		response.WriteError(w, r, h.logger, err)
		return
	}

	w.Header().Set("Location", "/api/v1/trips/"+t.ID)
	response.WriteJSON(w, http.StatusCreated, tripDTO)
}
