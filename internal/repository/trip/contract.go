package trip

import (
	"context"

	"github.com/swokeqq/tripngo.git/internal/infra/postgres/transaction"
)

type executorGetter interface {
	Executor(ctx context.Context) transaction.Querier
}
