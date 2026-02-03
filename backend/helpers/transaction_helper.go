package helpers

import (
	"context"

	"gorm.io/gorm"
)

// TransactionFunc represents a function that runs within a database transaction
type TransactionFunc func(ctx context.Context, tx *gorm.DB) error

// RunInTransaction executes a function within a database transaction
// It automatically handles commit and rollback based on the function's return value
func RunInTransaction(ctx context.Context, db *gorm.DB, fn TransactionFunc) error {
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-panic after rollback
		}
	}()

	if err := fn(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// RunInTransactionWithResult executes a function within a transaction and returns a result
func RunInTransactionWithResult[T any](ctx context.Context, db *gorm.DB, fn func(ctx context.Context, tx *gorm.DB) (T, error)) (T, error) {
	var zero T
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-panic after rollback
		}
	}()

	result, err := fn(ctx, tx)
	if err != nil {
		tx.Rollback()
		return zero, err
	}

	if err := tx.Commit().Error; err != nil {
		return zero, err
	}

	return result, nil
}
