package bot

import (
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/money"
)

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
		"• /total — итоги за текущий квартал (сумма и налог 6%)\n" +
		"• /help — подробная справка\n\n" +
		"Формат суммы: без знака минус, поддерживаются «1 234,56», «1234.56», «10р 50к». Деньги считаем детерминированно в копейках."
}

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

// BadAmountHintText returns a short hint for invalid /add amount input.
func BadAmountHintText() string {
	return "Не понял сумму. Примеры: 1000 | 1 234,56 | 10р 50к"
}

func UnknownCommandText() string {
	return "Неизвестная команда. Напишите /help"
}

func ErrorText() string {
	return "Ошибка при обработке команды. Попробуйте позже."
}

func AddSuccessText(amount int64, at time.Time, note string) string {
	// Deterministic template reply (no AI).
	var b strings.Builder
	b.WriteString("Добавлено: ")
	b.WriteString(money.FormatAmountShort(amount))
	b.WriteString("\nДата: ")
	b.WriteString(at.Format("2006-01-02"))
	if note != "" {
		b.WriteString("\nКомментарий: ")
		b.WriteString(note)
	}

	return b.String()
}

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
