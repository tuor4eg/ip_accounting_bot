package telegram

// User is a Telegram user or bot (minimal subset).
type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// Chat identifies the conversation the message belongs to.
type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title,omitempty"`
	Username string `json:"username,omitempty"`
}

// Message is a Telegram message (minimal subset for text echo).
type Message struct {
	MessageID int64  `json:"message_id"`
	Date      int64  `json:"date"`
	Text      string `json:"text,omitempty"`
	Chat      Chat   `json:"chat"`
	From      *User  `json:"from,omitempty"`
}

// Update is a Telegram Bot API update object (reduced to what we need now).
type Update struct {
	UpdateID      int64    `json:"update_id"`
	Message       *Message `json:"message,omitempty"`
	EditedMessage *Message `json:"edited_message,omitempty"`
}

// GetUpdatesParams are parameters for long polling.
type GetUpdatesParams struct {
	Offset         int64    `json:"offset,omitempty"`          // next update_id to receive
	Timeout        int      `json:"timeout,omitempty"`         // seconds to hold the long poll
	AllowedUpdates []string `json:"allowed_updates,omitempty"` // e.g. []{"message"}
}

// SendMessageParams are parameters for sending a text message.
type SendMessageParams struct {
	ChatID              int64  `json:"chat_id"`
	Text                string `json:"text"`
	ParseMode           string `json:"parse_mode,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ReplyToMessageID    int64  `json:"reply_to_message_id,omitempty"`
}
