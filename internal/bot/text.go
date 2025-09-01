package bot

import (
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/money"
)

// ------------------ COMMON MESSAGES ------------------
func UnknownCommandText() string {
	return "Неизвестная команда. Напишите /help"
}

func ErrorText() string {
	return "Ошибка при обработке команды. Попробуйте позже."
}

// ------------------ START MESSAGE ------------------

// StartText returns the greeting and quick usage guide for the bot.
// Text is static and transport-agnostic; actual sending is done by the router/runner.
func StartText() string {
	return "" +
		"Привет! Я помогу вести учёт доходов ИП (УСН 6%).\n\n" +
		"Команды:\n" +
		"• /add <сумма> [комментарий] — добавить поступление\n" +
		"  Примеры: /add 1000\n" +
		"           /add 1 234,56 заказ #42\n" +
		"           /add 10р 50к аванс\n" +
		"• /add_contrib <сумма> [комментарий] — добавить взнос\n" +
		"• /add_advance <сумма> [комментарий] — добавить авансовый платеж\n" +
		"• /undo — отменить последнее поступление за квартал\n" +
		"• /undo_contrib — отменить последний взнос\n" +
		"• /undo_advance — отменить последний авансовый платеж\n" +
		"• /total — итоги за текущий квартал (сумма и налог 6%)\n" +
		"• /help — подробная справка\n\n" +
		"Формат суммы: без знака минус, поддерживаются «1 234,56», «1234.56», «10р 50к»."
}

// ------------------ HELP MESSAGE ------------------

// HelpText returns a longer help message for users.
// Text is static and transport-agnostic.
func HelpText() string {
	return "" +
		"Справка\n" +
		"\n" +
		"Команды:\n" +
		"• /add <сумма> [комментарий]\n" +
		"  Добавляет поступление в базу. Сумма — без минуса, в рублях и копейках.\n" +
		"  Примеры:\n" +
		"   /add 1000\n" +
		"   /add 1 234,56 заказ #42\n" +
		"   /add 10р 50к аванс\n" +
		"\n" +
		"• /add_contrib <сумма> [комментарий]\n" +
		"  Добавляет взнос в базу. Сумма — аналогично /add.\n" +
		"\n" +
		"• /add_advance <сумма> [комментарий]\n" +
		"  Добавляет авансовый платеж в базу. Сумма — аналогично /add.\n" +
		"\n" +
		"• /undo\n" +
		"  Отменяет последнее поступление за квартал.\n" +
		"\n" +
		"• /undo_contrib\n" +
		"  Отменяет последний взнос.\n" +
		"\n" +
		"• /undo_advance\n" +
		"  Отменяет последний авансовый платеж.\n" +
		"\n" +
		"• /total\n" +
		"  Показывает сумму доходов и налог 6% за текущий квартал.\n" +
		"\n" +
		"• /start\n" +
		"  Краткая инструкция.\n" +
		"\n" +
		"Формат суммы:\n" +
		"  • Допускаются пробелы/точки/запятые как разделители тысяч.\n" +
		"  • Последняя точка или запятая — десятичный разделитель (до 2 знаков).\n" +
		"  • Понимает записи вида «10р 50к», «10 руб 50 коп».\n" +
		"  • Отрицательные значения не принимаются.\n" +
		"\n" +
		"Механика:\n" +
		"  • Налог рассчитывается как 6% от суммы квартала (округление вниз).\n" +
		"  • Квартал определяется по UTC датам (включительно).\n"
}

// ------------------ ADD MESSAGE ------------------

// BadAmountHintText returns a short hint for invalid /add amount input.
func BadAmountHintText() string {
	return "Не понял сумму. Примеры: 1000 | 1 234,56 | 10р 50к"
}

func AmountIsZeroText() string {
	return "Сумма не может быть 0"
}

func AddSuccessText(amount int64, at time.Time, note string) string {
	// Deterministic template reply (no AI).
	var b strings.Builder
	b.WriteString("✅ Добавлено поступление: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\nДата: ")
	b.WriteString(at.Format("2006-01-02"))
	if note != "" {
		b.WriteString("\nКомментарий: ")
		b.WriteString(note)
	}

	return b.String()
}

// ------------------ ADD CONTRIB MESSAGE ------------------

func AddContribSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Добавлен взнос: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\nДата: ")
	b.WriteString(at.Format("2006-01-02"))
	if note != "" {
		b.WriteString("\nКомментарий: ")
		b.WriteString(note)
	}
	b.WriteString(note)
	return b.String()
}

// ------------------ ADD ADVANCE MESSAGE ------------------

func AddAdvanceSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Добавлен авансовый платеж: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\nДата: ")
	b.WriteString(at.Format("2006-01-02"))
	if note != "" {
		b.WriteString("\nКомментарий: ")
		b.WriteString(note)
	}
	b.WriteString(note)
	return b.String()
}

// ------------------ TOTAL MESSAGE ------------------

func TotalText(sum int64, tax int64, qStart time.Time, qEnd time.Time) string {

	var b strings.Builder
	b.WriteString("Сумма за квартал: ")
	b.WriteString(qStart.Format("2006-01-02"))
	b.WriteString(" - ")
	b.WriteString(qEnd.Format("2006-01-02"))
	b.WriteString("\nСумма: ")
	b.WriteString(money.FormatAmountShort(sum))
	b.WriteString("\nНалог: ")
	b.WriteString(money.FormatAmountShort(tax))

	return b.String()
}

// ------------------ UNDO MESSAGE ------------------

func UndoSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Поступление отменено:\n")
	b.WriteString("\tСумма: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\n\tДата: ")
	b.WriteString(at.Format("2006-01-02"))
	if note != "" {
		b.WriteString("\n\tКомментарий: ")
		b.WriteString(note)
	}
	return b.String()
}

func UndoNoIncomeText() string {
	return "Нечего отменять. Нет поступлений за текущий квартал."
}

// ------------------ UNDO CONTRIB MESSAGE ------------------

func UndoContribSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Взнос отменен:\n")
	b.WriteString("\tСумма: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\n\tДата: ")
	b.WriteString(at.Format("2006-01-02"))
	if note != "" {
		b.WriteString("\n\tКомментарий: ")
		b.WriteString(note)
	}
	return b.String()
}

func UndoNoContribText() string {
	return "Нечего отменять. Нет взносов за текущий год."
}

// ------------------ UNDO ADVANCE MESSAGE ------------------

func UndoAdvanceSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Авансовый платеж отменен:\n")
	b.WriteString("\tСумма: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\n\tДата: ")
	b.WriteString(at.Format("2006-01-02"))
	if note != "" {
		b.WriteString("\n\tКомментарий: ")
		b.WriteString(note)
	}
	return b.String()
}

func UndoNoAdvanceText() string {
	return "Нечего отменять. Нет авансовых платежей за текущий год."
}
