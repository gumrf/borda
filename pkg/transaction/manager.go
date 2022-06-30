package transaction

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionManager interface {
	Run(ctx context.Context, fn func(ctx context.Context) error)
}

type transactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	// tm := new(transactionManager)
	return &transactionManager{pool: pool}
}

func (tm *transactionManager) Run(ctx context.Context, fn func(ctx context.Context) error) {
	// Extract transaction from context
	// Otherwise begin transaction
	// Inject transaction
	// Execute function
	// Checl errors
	// Commit if needed
}
