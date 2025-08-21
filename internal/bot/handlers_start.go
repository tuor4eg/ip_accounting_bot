package bot

import "context"

// HandleStart returns the response text for the /start command.
// Transport-agnostic; the router/runner is responsible for delivery.
func HandleStart(ctx context.Context) string {
	_ = ctx // reserved for future use (timeouts, locale, etc.)
	return StartText()
}
