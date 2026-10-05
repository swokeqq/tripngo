package handler

import (
	api "github.com/swokeqq/tripngo.git/internal/generated"
)

type Handler struct {
	api.Unimplemented
}

func New() *Handler {
	return &Handler{}
}
