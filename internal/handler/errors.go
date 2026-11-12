package handler

import "net/http"

type problemKind struct {
	code   string
	status int
	title  string
	typ    string
}

var (
	problemInvalidRequest = problemKind{
		typ:    "https://tripgo.example/problems/invalid-request",
		code:   "invalid_request",
		status: http.StatusBadRequest,
		title:  "Invalid request",
	}
	problemTripNotFound = problemKind{
		typ:    "https://tripgo.example/problems/trip-not-found",
		code:   "trip_not_found",
		status: http.StatusNotFound,
		title:  "Trip not found",
	}
	problemTripCompleted = problemKind{
		typ:    "https://tripgo.example/problems/trip-completed",
		code:   "trip_completed",
		status: http.StatusConflict,
		title:  "Trip completed",
	}
	problemDriverBusy = problemKind{
		typ:    "https://tripgo.example/problems/driver-busy",
		code:   "driver_busy",
		status: http.StatusConflict,
		title:  "Driver busy",
	}
	problemInternal = problemKind{
		typ:    "https://tripgo.example/problems/internal-error",
		code:   "internal_error",
		status: http.StatusInternalServerError,
		title:  "Internal Server error",
	}
)
