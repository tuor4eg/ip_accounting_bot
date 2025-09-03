package bot

import (
	"strconv"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/money"
)

// ------------------ COMMON MESSAGES ------------------
func UnknownCommandText() string {
	var b strings.Builder
	b.WriteString("❓ Неизвестная команда. Напишите /help для справки.")
	return b.String()
}

func ErrorText() string {
	var b strings.Builder
	b.WriteString("⚠️ Ошибка при обработке команды. Попробуйте позже.")
	return b.String()
}

// ------------------ START MESSAGE ------------------

// StartText returns the greeting and quick usage guide for the bot.
// Text is static and transport-agnostic; actual sending is done by the router/runner.
func StartText() string {
	var b strings.Builder
	b.WriteString("👋 Привет! Я помогу вести учёт доходов ИП (УСН 6%).\n\n")
	b.WriteString("📋 Основные команды:\n")
	b.WriteString("• /add [сумма] [комментарий] — добавить поступление\n")
	b.WriteString("  Примеры: /add 1000\n")
	b.WriteString("           /add 1 234,56 заказ #42\n")
	b.WriteString("           /add 10р 50к аванс\n")
	b.WriteString("• /add_contrib [сумма] [комментарий] — добавить взнос\n")
	b.WriteString("• /add_advance [сумма] [комментарий] — добавить авансовый платеж\n")
	b.WriteString("• /undo — отменить последнее поступление за квартал\n")
	b.WriteString("• /undo_contrib — отменить последний взнос\n")
	b.WriteString("• /undo_advance — отменить последний авансовый платеж\n")
	b.WriteString("• /total — итоги за текущий квартал (сумма и налог 6%)\n")
	b.WriteString("• /help — подробная справка\n\n")
	b.WriteString("💡 Формат суммы: без знака минус, поддерживаются «1 234,56», «1234.56», «10р 50к».")
	return b.String()
}

// ------------------ HELP MESSAGE ------------------

// HelpText returns a longer help message for users.
// Text is static and transport-agnostic.
func HelpText() string {
	var b strings.Builder
	b.WriteString("📚 Справка\n\n")
	b.WriteString("🔧 Команды:\n")
	b.WriteString("• /add [сумма] [комментарий]\n")
	b.WriteString("  Добавляет поступление в базу. Сумма — без минуса, в рублях и копейках.\n")
	b.WriteString("  Примеры:\n")
	b.WriteString("   /add 1000\n")
	b.WriteString("   /add 1 234,56 заказ #42\n")
	b.WriteString("   /add 10р 50к аванс\n\n")
	b.WriteString("• /add_contrib [сумма] [комментарий]\n")
	b.WriteString("  Добавляет взнос в базу. Сумма — аналогично /add.\n\n")
	b.WriteString("• /add_advance [сумма] [комментарий]\n")
	b.WriteString("  Добавляет авансовый платеж в базу. Сумма — аналогично /add.\n\n")
	b.WriteString("• /undo\n")
	b.WriteString("  Отменяет последнее поступление за квартал.\n\n")
	b.WriteString("• /undo_contrib\n")
	b.WriteString("  Отменяет последний взнос.\n\n")
	b.WriteString("• /undo_advance\n")
	b.WriteString("  Отменяет последний авансовый платеж.\n\n")
	b.WriteString("• /total\n")
	b.WriteString("  Показывает сумму доходов и налог 6% за текущий квартал.\n\n")
	b.WriteString("• /start\n")
	b.WriteString("  Краткая инструкция.\n\n")
	b.WriteString("💰 Формат суммы:\n")
	b.WriteString("  • Допускаются пробелы/точки/запятые как разделители тысяч.\n")
	b.WriteString("  • Последняя точка или запятая — десятичный разделитель (до 2 знаков).\n")
	b.WriteString("  • Понимает записи вида «10р 50к», «10 руб 50 коп».\n")
	b.WriteString("  • Отрицательные значения не принимаются.\n\n")
	b.WriteString("⚙️ Механика:\n")
	b.WriteString("  • Налог рассчитывается как 6% от суммы квартала (округление вниз).\n")
	b.WriteString("  • Квартал определяется по UTC датам (включительно).\n")
	return b.String()
}

// ------------------ ADD MESSAGE ------------------

// BadAmountHintText returns a short hint for invalid /add amount input.
func BadAmountHintText() string {
	var b strings.Builder
	b.WriteString("❌ Не понял сумму. Примеры: 1000 | 1 234,56 | 10р 50к")
	return b.String()
}

func AmountIsZeroText() string {
	var b strings.Builder
	b.WriteString("❌ Сумма не может быть 0")
	return b.String()
}

func AddSuccessText(amount int64, at time.Time, note string) string {
	// Deterministic template reply (no AI).
	var b strings.Builder
	b.WriteString("✅ Добавлено поступление: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\n📅 Дата: ")
	b.WriteString(at.Format("02.01.2006"))
	if note != "" {
		b.WriteString("\n💬 Комментарий: ")
		b.WriteString(note)
	}

	return b.String()
}

// ------------------ ADD CONTRIB MESSAGE ------------------

func AddContribSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Добавлен взнос: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\n📅 Дата: ")
	b.WriteString(at.Format("02.01.2006"))
	if note != "" {
		b.WriteString("\n💬 Комментарий: ")
		b.WriteString(note)
	}
	return b.String()
}

// ------------------ ADD ADVANCE MESSAGE ------------------

func AddAdvanceSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Добавлен авансовый платеж: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\n📅 Дата: ")
	b.WriteString(at.Format("02.01.2006"))
	if note != "" {
		b.WriteString("\n💬 Комментарий: ")
		b.WriteString(note)
	}
	return b.String()
}

// ------------------ TOTAL MESSAGE ------------------

func TotalText(
	sum int64, tax int64, qStart time.Time, qEnd time.Time,
	yearSum int64, yearTax int64,
	contribSum int64, advanceSum int64,
	year int, quarter int,
) string {
	var b strings.Builder

	// Quarter section
	b.WriteString("📅 <b>")
	b.WriteString(strconv.Itoa(quarter))
	b.WriteString(" квартал: ")
	b.WriteString(qStart.Format("02.01.2006"))
	b.WriteString(" - ")
	b.WriteString(qEnd.Format("02.01.2006"))
	b.WriteString("</b>")
	b.WriteString("\n")

	b.WriteString("💰 Поступления: ")
	b.WriteString(money.FormatAmountShort(sum))
	b.WriteString("\n")

	b.WriteString("🧾 Налог: ")
	b.WriteString(money.FormatAmountShort(tax))
	b.WriteString("\n\n")

	// Year section
	b.WriteString("📊 <b>Итого за ")
	b.WriteString(strconv.Itoa(year))
	b.WriteString(" год:</b>")
	b.WriteString("\n")

	b.WriteString("💰 Поступления: ")
	b.WriteString(money.FormatAmountShort(yearSum))
	b.WriteString("\n")

	b.WriteString("💳 Взносы: ")
	b.WriteString(money.FormatAmountShort(contribSum))
	b.WriteString("\n")

	b.WriteString("💸 Авансы: ")
	b.WriteString(money.FormatAmountShort(advanceSum))
	b.WriteString("\n")

	b.WriteString("🧾 Налог: ")
	b.WriteString(money.FormatAmountShort(yearTax))

	return b.String()
}

// ------------------ UNDO MESSAGE ------------------

func UndoSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Поступление отменено:\n")
	b.WriteString("💰 Сумма: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\n📅 Дата: ")
	b.WriteString(at.Format("02.01.2006"))
	if note != "" {
		b.WriteString("\n💬 Комментарий: ")
		b.WriteString(note)
	}
	return b.String()
}

func UndoNoIncomeText() string {
	var b strings.Builder
	b.WriteString("ℹ️ Нечего отменять. Нет поступлений за текущий квартал.")
	return b.String()
}

// ------------------ UNDO CONTRIB MESSAGE ------------------

func UndoContribSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Взнос отменен:\n")
	b.WriteString("💰 Сумма: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\n📅 Дата: ")
	b.WriteString(at.Format("02.01.2006"))
	if note != "" {
		b.WriteString("\n💬 Комментарий: ")
		b.WriteString(note)
	}
	return b.String()
}

func UndoNoContribText() string {
	var b strings.Builder
	b.WriteString("ℹ️ Нечего отменять. Нет взносов за текущий год.")
	return b.String()
}

// ------------------ UNDO ADVANCE MESSAGE ------------------

func UndoAdvanceSuccessText(amount int64, at time.Time, note string) string {
	var b strings.Builder
	b.WriteString("✅ Авансовый платеж отменен:\n")
	b.WriteString("💰 Сумма: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\n📅 Дата: ")
	b.WriteString(at.Format("02.01.2006"))
	if note != "" {
		b.WriteString("\n💬 Комментарий: ")
		b.WriteString(note)
	}
	return b.String()
}

func UndoNoAdvanceText() string {
	var b strings.Builder
	b.WriteString("ℹ️ Нечего отменять. Нет авансовых платежей за текущий год.")
	return b.String()
}
