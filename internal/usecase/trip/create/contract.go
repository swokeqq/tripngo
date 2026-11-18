//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package create

import (
	"context"

	"github.com/swokeqq/tripngo.git/internal/domain/trip"
	"github.com/swokeqq/tripngo.git/internal/model"
)

type (
	tripRepository interface {
		Save(ctx context.Context, t *trip.Trip) error
	}

	historyRepository interface {
		Add(ctx context.Context, history *model.TripStatusHistory) error
	}

	txManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}
)
