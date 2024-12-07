package dbutils

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

var (
	ErrCouldNotAcquireLock = errors.New("could not acquire database lock")
)

type Options struct {
	locks          []string
	timeoutSeconds uint8
}

// WithAdvisoryLock this option configures advisory locks to the given transaction.
func WithAdvisoryLock(key string) func(*Options) {
	return func(t *Options) {
		t.locks = append(t.locks, key)
	}
}

// WithLockTimeout this option configures the local lock timeout.
func WithLockTimeout(timeoutSeconds uint8) func(*Options) {
	return func(t *Options) {
		t.timeoutSeconds = timeoutSeconds
	}
}

// Transaction opens a transaction with the possibility of special options that are bound to it. The code that runs in
// the do parameter function is fully transactional with all its options.
func Transaction(
	ctx context.Context,
	db *pgx.Conn,
	do func(tx pgx.Tx) error,
	options ...func(*Options),
) error {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	if opts.timeoutSeconds > 0 {
		_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL lock_timeout = '%ds';", opts.timeoutSeconds))
		if err != nil {
			return err
		}
	}
	if err := NewPGXLocks(ctx, tx, opts.locks...); err != nil {
		if strings.Contains(err.Error(), "ERROR: canceling statement due to lock timeout (SQLSTATE 55P03)") {
			return ErrCouldNotAcquireLock
		}
		return err
	}
	if err := do(tx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("failed to roll back transaction: %w", err)
		}
		return fmt.Errorf("transaction rollback: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
