package app

import "errors"

var (
	ErrStoreNotSet                        = errors.New("store is not set")
	ErrIncomeUsecaseNotSet                = errors.New("income usecase is not set")
	ErrStoreDoesNotImplementIdentityStore = errors.New("store does not implement IdentityStore")
	ErrPaymentUsecaseNotSet               = errors.New("payment usecase is not set")
	ErrTotalUsecaseNotSet                 = errors.New("total usecase is not set")
)
