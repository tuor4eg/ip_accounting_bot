package bot_test

import (
	"context"
	"testing"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/service"
	"github.com/tuor4eg/ip_accounting_bot/internal/storage/memstore"
)

// Mock structures for tests
type mockPaymentService struct{}

func (m *mockPaymentService) AddPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, payoutType domain.PaymentType) error {
	return nil
}

func (m *mockPaymentService) UndoLastYear(ctx context.Context, userID int64, now time.Time, paymentType domain.PaymentType) (int64, time.Time, string, domain.PaymentType, bool, error) {
	return 0, time.Time{}, "", "", false, nil
}

type mockTotalService struct{}

func (m *mockTotalService) SumQuarter(ctx context.Context, userID int64, now time.Time) (domain.Totals, error) {
	return domain.Totals{}, nil
}

func (m *mockTotalService) SumYearToDate(ctx context.Context, userID int64, now time.Time) (domain.Totals, error) {
	return domain.Totals{}, nil
}

func fixedNow() time.Time {
	// inside quarter, UTC
	return time.Date(2025, 8, 10, 12, 0, 0, 0, time.UTC)
}

func TestHandleAdd_UsesTextTemplate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := memstore.NewStore()
	svc := service.NewIncomeService(store)

	// Create mock objects for Payment and Total
	mockPayment := &mockPaymentService{}
	mockTotal := &mockTotalService{}

	deps := &bot.BotDeps{
		Identities: store,
		Income:     svc,
		Payment:    mockPayment,
		Total:      mockTotal,
		Now:        fixedNow,
	}

	const transport = "telegram"
	const externalID = "42"
	const note = "тест"

	reply, err := bot.HandleAdd(ctx, deps, transport, externalID, "10р 50к "+note)
	if err != nil {
		t.Fatalf("HandleAdd error: %v", err)
	}

	// expected string — exactly through your template
	amount := int64(1050) // 10р 50к
	at := fixedNow()      // date in response is formatted as YYYY-MM-DD
	want := bot.AddSuccessText(amount, at, note)

	if reply != want {
		t.Fatalf("unexpected reply:\n--- got ---\n%s\n--- want ---\n%s", reply, want)
	}
}

func TestHandleUndo_NoIncome_UsesTextTemplate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := memstore.NewStore()
	svc := service.NewIncomeService(store)

	deps := &bot.BotDeps{
		Identities: store,
		Income:     svc, // implements UndoLastQuarter
		Payment:    &mockPaymentService{},
		Total:      &mockTotalService{},
		Now:        fixedNow,
	}

	const transport = "telegram"
	const externalID = "42"

	reply, err := bot.HandleUndo(ctx, deps, transport, externalID, "")
	if err != nil {
		t.Fatalf("HandleUndo error: %v", err)
	}

	want := bot.UndoNoIncomeText()
	if reply != want {
		t.Fatalf("unexpected reply:\n--- got ---\n%s\n--- want ---\n%s", reply, want)
	}
}

func TestHandleUndo_Success_UsesTextTemplate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := memstore.NewStore()
	svc := service.NewIncomeService(store)

	deps := &bot.BotDeps{
		Identities: store,
		Income:     svc,
		Payment:    &mockPaymentService{},
		Total:      &mockTotalService{},
		Now:        fixedNow,
	}

	const transport = "telegram"
	const externalID = "42"

	// first add income, then undo
	if _, err := bot.HandleAdd(ctx, deps, transport, externalID, "2 комментарий"); err != nil {
		t.Fatalf("HandleAdd error: %v", err)
	}

	reply, err := bot.HandleUndo(ctx, deps, transport, externalID, "")
	if err != nil {
		t.Fatalf("HandleUndo error: %v", err)
	}

	amount := int64(200) // 2 rubles = 200 kopecks
	at := fixedNow()
	note := "комментарий"
	want := bot.UndoSuccessText(amount, at, note)

	if reply != want {
		t.Fatalf("unexpected reply:\n--- got ---\n%s\n--- want ---\n%s", reply, want)
	}
}
