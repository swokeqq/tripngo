package server

import (
	"errors"
	"log/slog"
	"net/http"

	tripdomain "github.com/swokeqq/tripngo.git/internal/domain/trip"
	api "github.com/swokeqq/tripngo.git/internal/generated"
	"github.com/swokeqq/tripngo.git/internal/http/health"
	"github.com/swokeqq/tripngo.git/internal/http/ready"
	"github.com/swokeqq/tripngo.git/internal/http/response"
	"github.com/swokeqq/tripngo.git/internal/http/trip/create"
	"github.com/swokeqq/tripngo.git/internal/http/trip/finish"
	"github.com/swokeqq/tripngo.git/internal/http/trip/get"
)

type Server struct {
	api.Unimplemented
	c      *create.Handler
	g      *get.Handler
	f      *finish.Handler
	rd     *ready.Handler
	logger *slog.Logger
}

func New(create *create.Handler, get *get.Handler,
	finish *finish.Handler, ready *ready.Handler, logger *slog.Logger,
) *Server {
	return &Server{c: create, g: get, f: finish, rd: ready, logger: logger}
}

func (s *Server) GetTrip(w http.ResponseWriter, r *http.Request, id api.TripId) {
	s.g.GetTrip(w, r, id)
}

func (s *Server) FinishTrip(w http.ResponseWriter, r *http.Request, id api.TripId) {
	s.f.FinishTrip(w, r, id)
}

func (s *Server) CreateTrip(w http.ResponseWriter, r *http.Request, params api.CreateTripParams) {
	s.c.CreateTrip(w, r, params)
}

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	health.Health(w, r)
}

func (s *Server) Ready(w http.ResponseWriter, r *http.Request) {
	s.rd.Ready(w, r)
}

func (s *Server) ParamErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var invalidParam *api.InvalidParamFormatError
	if errors.As(err, &invalidParam) {
		response.WriteError(
			w, r, s.logger,
			&tripdomain.ValidationError{Field: invalidParam.ParamName, Reason: "invalid format"},
		)
		return
	}

	response.WriteError(w, r, s.logger,
		&tripdomain.ValidationError{Field: "body", Reason: err.Error()})
}

var _ api.ServerInterface = (*Server)(nil)
