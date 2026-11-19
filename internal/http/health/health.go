package health

import (
	"net/http"

	api "github.com/swokeqq/tripngo.git/internal/generated"
	"github.com/swokeqq/tripngo.git/internal/http/response"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	response.WriteJSON(w, http.StatusOK, api.HealthResponse{Status: api.Ok})
}
