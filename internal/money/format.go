package money

import (
	"fmt"
)

// Currency symbols and names
const (
	RubleSymbol = "₽"
	KopekSymbol = "коп"
)

func FormatAmount(amount int64) (string, error) {
	rubs := amount / 100
	kopeks := amount % 100

	return fmt.Sprintf("%d %s %02d %s", rubs, RubleSymbol, kopeks, KopekSymbol), nil
}

// FormatAmountShort formats amount as "123.45₽" (compact format)
func FormatAmountShort(amount int64) string {
	rubs := amount / 100
	kopeks := amount % 100
	return fmt.Sprintf("%d.%02d%s", rubs, kopeks, RubleSymbol)
}

// GetAmountParts returns rubles and kopecks separately
func GetAmountParts(amount int64) (rubles, kopecks int64) {
	return amount / 100, amount % 100
}
