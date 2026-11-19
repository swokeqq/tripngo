package trip

import (
	"time"

	"github.com/Masterminds/squirrel"
)

const pgUniqueViolation = "23505"

type HistoryRepository struct {
	db           executorGetter
	queryTimeout time.Duration
	qb           squirrel.StatementBuilderType
}

func NewHistoryRepository(db executorGetter, queryTimeout time.Duration) *HistoryRepository {
	return &HistoryRepository{
		db:           db,
		queryTimeout: queryTimeout,
		qb:           squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

type Repository struct {
	db           executorGetter
	queryTimeout time.Duration
	qb           squirrel.StatementBuilderType
}

func NewRepository(db executorGetter, queryTimeout time.Duration) *Repository {
	return &Repository{
		db:           db,
		queryTimeout: queryTimeout,
		qb:           squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
