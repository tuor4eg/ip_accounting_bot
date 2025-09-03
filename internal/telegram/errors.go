package telegram

import "fmt"

// APIError represents a Telegram API error
type APIError struct {
	Status      int
	Description string
}

func (e APIError) Error() string {
	if e.Status != 0 && e.Description != "" {
		return fmt.Sprintf("Telegram API: Status=%d, %s", e.Status, e.Description)
	}

	if e.Description != "" {
		return fmt.Sprintf("Telegram API: %s", e.Description)
	}

	return "Telegram API error"
}
