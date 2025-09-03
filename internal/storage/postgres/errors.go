package postgres

import "errors"

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
