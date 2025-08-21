package migrations

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UnlockFunc releases the advisory lock acquired by AcquireAdvisoryLock.
type UnlockFunc func(ctx context.Context) error

// AcquireAdvisoryLock takes an application-scoped lock name, blocks until the lock
// is acquired on a dedicated connection from the pool, and returns a function to
// release it. The same connection is held until Unlock is called.
//
// Usage:
//   unlock, err := AcquireAdvisoryLock(ctx, pool, "ip_accounting_bot:migrations")
//   if err != nil { ... }
//   defer unlock(context.Background())

func AcquireAdvisoryLock(ctx context.Context, pool *pgxpool.Pool, name string) (UnlockFunc, error) {
	if pool == nil {
		return nil, fmt.Errorf("pool is nil")
	}

	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	// Stable 64-bit key from name (negative values are fine for Postgres BIGINT).
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(name))
	lockID := int64(hash.Sum64())

	conn, err := pool.Acquire(ctx)

	if err != nil {
		return nil, fmt.Errorf("acquire connection: %w", err)
	}

	if _, err := conn.Exec(ctx, `SELECT pg_advisory_lock($1)`, lockID); err != nil {
		conn.Release()

		return nil, fmt.Errorf("acquire advisory lock: pg_advisory_lock(%d): %w", lockID, err)
	}

	unlock := func(ctx context.Context) error {
		defer conn.Release()

		if ctx == nil {
			ctx = context.Background()
		}

		if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) <= 0 {
			var cancel context.CancelFunc

			ctx, cancel = context.WithTimeout(ctx, 5*time.Second)

			defer cancel()
		}

		if _, err := conn.Exec(ctx, `SELECT pg_advisory_unlock($1)`, lockID); err != nil {
			return fmt.Errorf("release advisory lock: pg_advisory_unlock(%d): %w", lockID, err)
		}

		return nil
	}

	return unlock, nil
}
