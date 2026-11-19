//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package get

import (
	"context"

	"github.com/swokeqq/tripngo.git/internal/domain/trip"
)

type usecase interface {
	GetByID(ctx context.Context, id string) (*trip.Trip, error)
}
