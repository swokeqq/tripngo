package handler

import (
	"encoding/json"
	"log"
	"net/http"

	api "github.com/swokeqq/tripngo.git/internal/generated"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := api.HealthResponse{
		Status: api.Ok,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
