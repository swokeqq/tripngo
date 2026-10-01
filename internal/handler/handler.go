package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	api "github.com/swokeqq/tripngo.git/internal/generated"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	db Pinger
}

func NewHandler(db Pinger) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, api.HealthResponse{Status: "ok"})
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.db == nil || h.db.Ping(r.Context()) != nil {
		respondJSON(w, http.StatusServiceUnavailable, api.HealthResponse{Status: "unavailable"})
		return
	}

	respondJSON(w, http.StatusOK, api.HealthResponse{Status: "ok"})
}

func (h *Handler) CreateTrip(w http.ResponseWriter, r *http.Request, tripParams api.CreateTripParams) {
	respondJSON(w, http.StatusServiceUnavailable, api.HealthResponse{Status: "not implemented yet"})
}

func (h *Handler) GetTrip(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusServiceUnavailable, api.HealthResponse{Status: "not implemented yet"})
}

func (h *Handler) FinishTrip(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusServiceUnavailable, api.HealthResponse{Status: "not implemented yet"})
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
