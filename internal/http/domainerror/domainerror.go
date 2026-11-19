package domainerror

import (
	"errors"
	"net/http"

	tripdomain "github.com/swokeqq/tripngo.git/internal/domain/trip"
	api "github.com/swokeqq/tripngo.git/internal/generated"
)

type problemKind struct {
	code   string
	status int32
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

func Problem(err error) (status int, body api.Problem) {
	var validationErr *tripdomain.ValidationError
	var p api.Problem
	if errors.As(err, &validationErr) {
		p = toProblem(problemInvalidRequest, validationErr)
		return int(p.Status), p
	}
	if errors.Is(err, tripdomain.ErrTripNotFound) {
		p = toProblem(problemTripNotFound, tripdomain.ErrTripNotFound)
		return int(p.Status), p
	}
	if errors.Is(err, tripdomain.ErrDriverBusy) {
		p = toProblem(problemDriverBusy, tripdomain.ErrDriverBusy)
		return int(p.Status), p
	}
	if errors.Is(err, tripdomain.ErrTripCompleted) {
		p = toProblem(problemTripCompleted, tripdomain.ErrTripCompleted)
		return int(p.Status), p
	}

	p = toProblem(problemInternal, errors.New("internal server error"))
	return int(p.Status), p
}

func toProblem(problem problemKind, err error) api.Problem {
	detail := err.Error()
	return api.Problem{
		Code:   problem.code,
		Detail: &detail,
		Status: problem.status,
		Title:  problem.title,
		Type:   problem.typ,
	}
}
