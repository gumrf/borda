package transaction

// type key struct{}

// func injectTransaction(ctx context.Context, tx pgx.Tx) context.Context {
// 	return context.WithValue(ctx, key{}, tx)
// }

// func ExtractTransaction(ctx context.Context) (pgx.Tx, bool) {
// 	if tx, ok := ctx.Value(key{}).(pgx.Tx); ok {
// 		return tx, true
// 	}
// 	return nil, false
// }

// func Transactional(ctx context.Context, db *pgxpool.Pool, fn func(ctx context.Context) error, isDone ...bool) error {
// 	var tx pgx.Tx
// 	var err error
// 	var ok bool

// 	tx, ok = ExtractTransaction(ctx)

// 	if !ok {
// 		tx, err = db.Begin(ctx)
// 		if err != nil {
// 			return fmt.Errorf("can't begin transaction: %w", err)
// 		}

// 		ctx = injectTransaction(ctx, tx)
// 	}

// 	if err := fn(ctx); err != nil {
// 		// TODO: What error should be returned???
// 		rollbackError := tx.Rollback(ctx)
// 		if rollbackError != nil {
// 			// if errors.Is(rollbackError, pgx.ErrTxClosed) {
// 			// 	return nil
// 			// }
// 			log.Printf("can't rollback transaction: %v", rollbackError)
// 		}
// 		return err
// 	}

// 	if len(isDone) > 0 && isDone[0] == true { //nolint
// 		fmt.Println("iam here")
// 		commitError := tx.Commit(ctx)
// 		if commitError != nil {
// 			log.Printf("can't commit transaction: %v", commitError)
// 			return commitError
// 		}
// 	}

// 	return nil
// }
