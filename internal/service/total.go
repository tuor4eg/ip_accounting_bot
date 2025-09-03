// internal/service/total.go
package service

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/period"
	"github.com/tuor4eg/ip_accounting_bot/internal/tax"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

// NewTotalService wires functional dependencies (no direct store coupling).
func NewTotalService(
	getUserScheme func(ctx context.Context, userID int64) (domain.TaxScheme, error),
	sumIncomes func(ctx context.Context, userID int64, from, to time.Time) (int64, error),
	sumPayments func(ctx context.Context, userID int64, from, to time.Time) (int64, int64, error),
	provider tax.Provider,
) *TotalService {
	return &TotalService{
		getUserScheme: getUserScheme,
		sumIncomes:    sumIncomes,
		sumPayments:   sumPayments,
		provider:      provider,
	}
}

func (s *TotalService) SumQuarter(ctx context.Context, userID int64, now time.Time) (domain.Totals, error) {
	return SumQuarter(
		ctx,
		s.getUserScheme,
		s.sumIncomes,
		s.sumPayments,
		s.provider,
		userID,
		now,
	)
}

func (s *TotalService) SumYearToDate(ctx context.Context, userID int64, now time.Time) (domain.Totals, error) {
	return SumYearToDate(
		ctx,
		s.getUserScheme,
		s.sumIncomes,
		s.sumPayments,
		s.provider,
		userID,
		now,
	)
}

// SumQuarter aggregates incomes and payments for the quarter that contains ref,
// selects tax policy for the user's scheme at the quarter end, and computes totals.
//   - 1% annual extra is NOT included here.
func SumQuarter(
	ctx context.Context,
	getUserScheme func(ctx context.Context, userID int64) (domain.TaxScheme, error),
	sumIncomes func(ctx context.Context, userID int64, from, to time.Time) (int64, error),
	sumPayments func(ctx context.Context, userID int64, from, to time.Time) (int64, int64, error), // (contrib, advance)
	provider tax.Provider,
	userID int64,
	ref time.Time,
) (domain.Totals, error) {
	const op = "service.total.SumQuarter"

	// Inclusive quarter bounds; storage layer trims to UTC DATE on write/read.
	from, to := period.QuarterBounds(ref)

	scheme, err := getUserScheme(ctx, userID)
	if err != nil {
		return domain.Totals{}, validate.Wrap(op, err)
	}

	policy, err := provider.ForDate(scheme, to)
	if err != nil {
		return domain.Totals{}, validate.Wrap(op, err)
	}

	incomeSum, err := sumIncomes(ctx, userID, from, to)
	if err != nil {
		return domain.Totals{}, validate.Wrap(op, err)
	}

	contribSum, advanceSum, err := sumPayments(ctx, userID, from, to)
	if err != nil {
		return domain.Totals{}, validate.Wrap(op, err)
	}

	// Deterministic integer math (kopecks only).

	taxAmount := incomeSum * policy.BaseRateBP / domain.BpDen

	contribApplied := min(contribSum, taxAmount)

	due := max(taxAmount-contribApplied-advanceSum, 0)

	return domain.Totals{
		From:           from,
		To:             to,
		IncomeSum:      incomeSum,
		Tax:            taxAmount,
		ContribSum:     contribSum,
		AdvanceSum:     advanceSum,
		ContribApplied: contribApplied,
		Due:            due,
	}, nil
}

func SumYearToDate(
	ctx context.Context,
	getUserScheme func(ctx context.Context, userID int64) (domain.TaxScheme, error),
	sumIncomes func(ctx context.Context, userID int64, from, to time.Time) (int64, error),
	sumPayments func(ctx context.Context, userID int64, from, to time.Time) (int64, int64, error), // (contrib, advance)
	provider tax.Provider,
	userID int64,
	ref time.Time,
) (domain.Totals, error) {
	const op = "service.total.SumYearToDate"

	from, to := period.YearBounds(ref)

	scheme, err := getUserScheme(ctx, userID)
	if err != nil {
		return domain.Totals{}, validate.Wrap(op, err)
	}

	policy, err := provider.ForDate(scheme, to)
	if err != nil {
		return domain.Totals{}, validate.Wrap(op, err)
	}

	incomeSum, err := sumIncomes(ctx, userID, from, to)
	if err != nil {
		return domain.Totals{}, validate.Wrap(op, err)
	}

	contribSum, advanceSum, err := sumPayments(ctx, userID, from, to)
	if err != nil {
		return domain.Totals{}, validate.Wrap(op, err)
	}

	taxAmount := incomeSum * policy.BaseRateBP / domain.BpDen

	contribApplied := min(contribSum, taxAmount)

	due := max(taxAmount-contribApplied-advanceSum, 0)

	return domain.Totals{
		From:           from,
		To:             to,
		IncomeSum:      incomeSum,
		Tax:            taxAmount,
		ContribSum:     contribSum,
		AdvanceSum:     advanceSum,
		ContribApplied: contribApplied,
		Due:            due,
	}, nil
}
