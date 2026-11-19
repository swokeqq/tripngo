package ready

import (
	"context"
	"net/http"
	"time"

	api "github.com/swokeqq/tripngo.git/internal/generated"
	"github.com/swokeqq/tripngo.git/internal/http/response"
)

type Handler struct {
	pinger       pinger
	queryTimeout time.Duration
}

func New(pinger pinger, queryTimeout time.Duration) *Handler {
	return &Handler{pinger: pinger, queryTimeout: queryTimeout}
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if err := h.pinger.Ping(ctx); err != nil {
		response.WriteJSON(w, http.StatusServiceUnavailable, api.HealthResponse{Status: api.Unavailable})
		return
	}

	response.WriteJSON(w, http.StatusOK, api.HealthResponse{Status: api.Ok})
}
