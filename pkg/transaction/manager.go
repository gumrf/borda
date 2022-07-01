package transaction

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Manager interface {
	Run(ctx context.Context, fn func(ctx context.Context) error, isDone ...bool) error
	GetTransaction(ctx context.Context) (pgx.Tx, bool)
}

type manager struct {
	pool *pgxpool.Pool
}

func NewManager(pool *pgxpool.Pool) Manager {
	return &manager{pool: pool}
}

type key struct{}

func (m *manager) GetTransaction(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(key{}).(pgx.Tx)
	return tx, ok
}

func injectTransaction(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, key{}, tx)
}

func (m *manager) Run(ctx context.Context, fn func(ctx context.Context) error, isDone ...bool) error {
	var tx pgx.Tx
	var err error
	var ok bool

	tx, ok = m.GetTransaction(ctx)

	if !ok {
		tx, err = m.pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("can't begin transaction: %w", err)
		}

		ctx = injectTransaction(ctx, tx)
	}

	if err := fn(ctx); err != nil {
		// TODO: What error should be returned???
		tx.Rollback(ctx)
		// if rollbackError != nil {
		// 	if errors.Is(rollbackError, pgx.ErrTxClosed) {
		// 		log.Printf("transaction closed: %v", rollbackError)
		// 	}

		// 	log.Printf("can't rollback transaction: %v", rollbackError)

		// 	// return rollbackError
		// }
		return err
	}

	if len(isDone) > 0 && isDone[0] == true { //nolint
		commitError := tx.Commit(ctx)
		if commitError != nil {
			log.Printf("can't commit transaction: %v", commitError)
			return commitError
		}
	}

	return nil
}
