package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tuor4eg/ip_accounting_bot/config"
	"github.com/tuor4eg/ip_accounting_bot/internal/logging"
	"github.com/tuor4eg/ip_accounting_bot/migrations"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Logs from env: LOG_LEVEL, LOG_FORMAT
	logging.InitFromEnv(cfg.LogLevel, cfg.LogFormat)

	// Flags
	var (
		dir     = flag.String("dir", "sql", "migrations directory inside embedded FS")
		timeout = flag.Duration("timeout", 5*time.Minute, "overall timeout for applying migrations")
	)
	flag.Parse()

	// Context with overall timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Connect pool
	cfgDb, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to parse DSN", "err", err)
		os.Exit(1)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfgDb)
	if err != nil {
		slog.Error("failed to create pool", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Quick connection check
	if err := pool.Ping(ctx); err != nil {
		slog.Error("database ping failed", "err", err)
		os.Exit(1)
	}

	// Apply missing migrations
	applied, err := migrations.ApplyUp(ctx, pool, migrations.FS, *dir)
	if err != nil {
		slog.Error("apply up migrations failed", "err", err)
		os.Exit(1)
	}

	slog.Info("migrations applied", "count", applied, "dir", *dir)
	fmt.Printf("OK: applied %d migration(s)\n", applied)
}
