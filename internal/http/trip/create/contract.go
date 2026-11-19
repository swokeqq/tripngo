//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package create

import (
	"context"

	"github.com/swokeqq/tripngo.git/internal/domain/trip"
	createuc "github.com/swokeqq/tripngo.git/internal/usecase/trip/create"
)

type usecase interface {
	Create(ctx context.Context, in createuc.Input) (*trip.Trip, error)
}
