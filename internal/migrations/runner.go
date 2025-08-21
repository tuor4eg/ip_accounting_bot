package migrations

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var reMigration = regexp.MustCompile(`^(?P<ver>\d{4,})_(?P<name>[a-z0-9_]+?)(?P<conc>_concurrently)?\.up\.sql$`)

type Migration struct {
	Version    int64
	Concurrent bool
	Path       string
	Name       string
}

func EnsureMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("pool is nil")
	}

	tx, err := pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func AppliedVersions(ctx context.Context, pool *pgxpool.Pool) ([]int64, error) {
	if pool == nil {
		return nil, fmt.Errorf("pool is nil")
	}

	rows, err := pool.Query(ctx, `
		SELECT version FROM schema_migrations
		ORDER BY version ASC
	`)

	if err != nil {
		return nil, fmt.Errorf("query schema_migrations: %w", err)
	}

	defer rows.Close()

	var vs []int64

	for rows.Next() {
		var v int64

		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("scan version: %w", err)
		}

		vs = append(vs, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return vs, nil
}

func ListUpMigrations(fsys fs.FS, dir string) ([]Migration, error) {
	if fsys == nil {
		return nil, fmt.Errorf("fsys is nil")
	}

	entries, err := fs.ReadDir(fsys, dir)

	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	var out []Migration
	seen := make(map[int64]string)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		v, conc, ok := parseMigrationName(name)

		if !ok {
			continue
		}

		if prev, dup := seen[v]; dup {
			return nil, fmt.Errorf("duplicate migration version %d: %q and %q", v, prev, name)
		}

		seen[v] = name

		out = append(out, Migration{
			Version:    v,
			Concurrent: conc,
			Path:       filepath.ToSlash(filepath.Join(dir, name)),
			Name:       name,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Version < out[j].Version
	})

	return out, nil
}

func ApplyUp(ctx context.Context, pool *pgxpool.Pool, fsys fs.FS, dir string) (int, error) {
	if pool == nil {
		return 0, fmt.Errorf("apply up: nil pool")
	}
	if fsys == nil {
		return 0, fmt.Errorf("apply up: nil fs")
	}

	// 1) Serialization of migrations (global lock).
	unlock, err := AcquireAdvisoryLock(ctx, pool, "ip_accounting_bot:migrations")
	if err != nil {
		return 0, err
	}
	defer func() { _ = unlock(context.Background()) }()

	// 2) Ensure that there is a service table.
	if err := EnsureMigrationsTable(ctx, pool); err != nil {
		return 0, err
	}

	// 3) Collect the list of available up-migrations and already applied versions.
	all, err := ListUpMigrations(fsys, dir)
	if err != nil {
		return 0, err
	}
	applied, err := AppliedVersions(ctx, pool)
	if err != nil {
		return 0, err
	}
	appliedSet := make(map[int64]struct{}, len(applied))
	for _, v := range applied {
		appliedSet[v] = struct{}{}
	}

	// 4) Apply the missing ones.
	appliedCount := 0
	for _, m := range all {
		if _, ok := appliedSet[m.Version]; ok {
			continue // already applied
		}

		sqlBytes, rerr := fs.ReadFile(fsys, m.Path)
		if rerr != nil {
			return appliedCount, fmt.Errorf("apply up: read %q: %w", m.Path, rerr)
		}
		sql := strings.TrimSpace(string(sqlBytes))

		if m.Concurrent {
			// 4a) Outside the transaction (for ... CONCURRENTLY).
			if _, err := pool.Exec(ctx, sql); err != nil {
				return appliedCount, fmt.Errorf("apply up: exec %q (v=%d): %w", m.Name, m.Version, err)
			}
			if _, err := pool.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES($1) ON CONFLICT DO NOTHING`, m.Version); err != nil {
				return appliedCount, fmt.Errorf("apply up: record version %d: %w", m.Version, err)
			}
		} else {
			// 4b) In the transaction.
			tx, err := pool.Begin(ctx)
			if err != nil {
				return appliedCount, fmt.Errorf("apply up: begin tx for %q (v=%d): %w", m.Name, m.Version, err)
			}
			// just in case, if something goes wrong
			defer func() { _ = tx.Rollback(ctx) }()

			if _, err := tx.Exec(ctx, sql); err != nil {
				_ = tx.Rollback(ctx)
				return appliedCount, fmt.Errorf("apply up: exec %q (v=%d): %w", m.Name, m.Version, err)
			}
			if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES($1)`, m.Version); err != nil {
				_ = tx.Rollback(ctx)
				return appliedCount, fmt.Errorf("apply up: record version %d: %w", m.Version, err)
			}
			if err := tx.Commit(ctx); err != nil {
				return appliedCount, fmt.Errorf("apply up: commit %q (v=%d): %w", m.Name, m.Version, err)
			}
		}

		appliedCount++
	}

	return appliedCount, nil
}

// parseMigrationName extracts version and "concurrent" flag from a migration file name.
// Valid examples:
//
//	0001_init.up.sql                        -> version=1,  concurrent=false
//	0002_big_index_concurrently.up.sql      -> version=2,  concurrent=true
func parseMigrationName(name string) (version int64, concurrent bool, ok bool) {
	base := filepath.Base(name)
	lower := strings.ToLower(base)

	m := reMigration.FindStringSubmatch(lower)
	if m == nil {
		return 0, false, false
	}

	verStr := m[reMigration.SubexpIndex("ver")]
	conc := m[reMigration.SubexpIndex("conc")] != ""

	v, err := strconv.ParseInt(verStr, 10, 64)
	if err != nil || v <= 0 {
		return 0, false, false
	}
	return v, conc, true
}
