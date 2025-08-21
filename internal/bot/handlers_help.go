package bot

import "context"

// HandleHelp returns the response text for the /help command.
// Transport-agnostic; the router/runner is responsible for delivery.
func HandleHelp(ctx context.Context) string {
	_ = ctx
	return HelpText()
}
