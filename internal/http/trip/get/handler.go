package get

import (
	"log/slog"
	"net/http"

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

func (h *Handler) GetTrip(w http.ResponseWriter, r *http.Request, tripId api.TripId) {
	t, err := h.uc.GetByID(r.Context(), tripId.String())
	if err != nil {
		response.WriteError(w, r, h.logger, err)
		return
	}

	tripDTO, err := response.ToAPITrip(t)
	if err != nil {
		response.WriteError(w, r, h.logger, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, tripDTO)
}
