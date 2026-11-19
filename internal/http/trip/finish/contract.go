//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package finish

import (
	"context"

	"github.com/swokeqq/tripngo.git/internal/domain/trip"
)

type usecase interface {
	Finish(ctx context.Context, tripID string) (*trip.Trip, error)
}
