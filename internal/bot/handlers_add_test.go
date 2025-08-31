package bot_test

import (
	"context"
	"testing"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
	"github.com/tuor4eg/ip_accounting_bot/internal/service"
	"github.com/tuor4eg/ip_accounting_bot/internal/storage/memstore"
)

func fixedNow() time.Time {
	// внутри квартала, UTC
	return time.Date(2025, 8, 10, 12, 0, 0, 0, time.UTC)
}

func TestHandleAdd_UsesTextTemplate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := memstore.NewStore()
	svc := service.NewIncomeService(store)

	deps := bot.AddDeps{
		Identities: store,
		Income:     svc,
		Now:        fixedNow,
	}

	const transport = "telegram"
	const externalID = "42"
	const note = "тест"

	reply, err := bot.HandleAdd(ctx, deps, transport, externalID, "10р 50к "+note)
	if err != nil {
		t.Fatalf("HandleAdd error: %v", err)
	}

	// ожидаемая строка — ровно через ваш шаблон
	amount := int64(1050) // 10р 50к
	at := fixedNow()      // дата в ответе форматируется как YYYY-MM-DD
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

	deps := bot.AddDeps{
		Identities: store,
		Income:     svc, // реализует UndoLastQuarter
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

	deps := bot.AddDeps{
		Identities: store,
		Income:     svc,
		Now:        fixedNow,
	}

	const transport = "telegram"
	const externalID = "42"

	// сначала добавим доход, затем отменим
	if _, err := bot.HandleAdd(ctx, deps, transport, externalID, "2 комментарий"); err != nil {
		t.Fatalf("HandleAdd error: %v", err)
	}

	reply, err := bot.HandleUndo(ctx, deps, transport, externalID, "")
	if err != nil {
		t.Fatalf("HandleUndo error: %v", err)
	}

	amount := int64(200) // 2 рубля = 200 коп
	at := fixedNow()
	note := "комментарий"
	want := bot.UndoSuccessText(amount, at, note)

	if reply != want {
		t.Fatalf("unexpected reply:\n--- got ---\n%s\n--- want ---\n%s", reply, want)
	}
}
