package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionManager manages database transactions
type TransactionManager interface {
	// WithTransaction executes a function within a transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error)
}

// PostgresTransactionManager implements TransactionManager for PostgreSQL
type PostgresTransactionManager struct {
	pool *pgxpool.Pool
}

// NewPostgresTransactionManager creates a new PostgreSQL transaction manager
func NewPostgresTransactionManager(pool *pgxpool.Pool) *PostgresTransactionManager {
	return &PostgresTransactionManager{
		pool: pool,
	}
}

// WithTransaction executes a function within a transaction
func (tm *PostgresTransactionManager) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) (interface{}, error),
) (interface{}, error) {
	// Begin transaction
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back on panic
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p) // Re-throw panic
		}
	}()

	// Add transaction to context
	txCtx := context.WithValue(ctx, txKey, tx)

	// Execute function
	result, err := fn(txCtx)
	if err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return nil, fmt.Errorf("tx error: %w, rollback error: %v", err, rbErr)
		}
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// Context key for transaction
type contextKey string

const txKey contextKey = "tx"

// GetTx retrieves the transaction from context
func GetTx(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		return nil
	}
	return tx
}
