package dbutils

import (
	"context"
	"fmt"

	"hash/fnv"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/4ND3R50N/go-tools/filter"
)

type AcquiredLocks struct {
	IDs []string
}

type PGXInterface interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// NewPGXLocks acquire pg_advisory_xact_lock locks for each lockID given. Duplicates in lockIDs getting filtered out.
func NewPGXLocks(ctx context.Context, db PGXInterface, lockIDs ...string) error {
	lockIDs = filter.Distinct(lockIDs)
	for _, id := range lockIDs {
		resourceHash := fnv.New64()
		_, _ = resourceHash.Write([]byte(id))
		_, err := db.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", int64(resourceHash.Sum64()))
		if err != nil {
			return fmt.Errorf("could not acquire database advisory lock: %w", err)
		}
	}
	return nil
}
