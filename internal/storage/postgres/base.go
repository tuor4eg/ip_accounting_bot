package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

type Store struct {
	Pool *pgxpool.Pool
}

const (
	PingTimeOutSec = 5
)

var (
	ErrEmptyDSN    = errors.New("dsn is empty")
	ErrEmptyPool   = errors.New("pool is empty")
	ErrParseConfig = errors.New("parse config error")
	ErrPoolCreate  = errors.New("pool create error")
	ErrPing        = errors.New("ping error")
	ErrBeginTx     = errors.New("begin tx error")
	ErrTx          = errors.New("tx error")
	ErrTxCommit    = errors.New("tx commit error")
)

// Open initializes a pgx pool from DSN and verifies the connection with Ping.
// DSN example: postgres://user:pass@host:5432/dbname?sslmode=disable
func Open(ctx context.Context, dsn string) (*Store, error) {
	const op = "postgres.Open"

	dsn = strings.TrimSpace(dsn)

	if dsn == "" {
		return nil, validate.Wrap(op, ErrEmptyDSN)
	}

	cfg, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, validate.Wrap(op, ErrParseConfig)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)

	if err != nil {
		return nil, validate.Wrap(op, ErrPoolCreate)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, validate.Wrap(op, ErrPing)
	}

	return &Store{Pool: pool}, nil
}

func (s *Store) Close(ctx context.Context) error {
	if s == nil || s.Pool == nil {
		return nil
	}

	s.Pool.Close()
	return nil
}

func (s *Store) WithTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	const op = "postgres.WithTx"

	if s == nil || s.Pool == nil {
		return validate.Wrap(op, ErrEmptyPool)
	}

	tx, err := s.Pool.Begin(ctx)

	if err != nil {
		return validate.Wrap(op, ErrBeginTx)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback(ctx)

		return validate.Wrap(op, ErrTx)
	}

	if err := tx.Commit(ctx); err != nil {
		return validate.Wrap(op, ErrTxCommit)
	}

	return nil
}
