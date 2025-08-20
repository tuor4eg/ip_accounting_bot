package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GetUpdates calls Telegram getUpdates with given params and returns a batch of updates.
// Note: allowed_updates must be JSON-encoded when sent via query/form (url.Values).
func (c *Client) GetUpdates(ctx context.Context, p GetUpdatesParams) ([]Update, error) {
	q := url.Values{}

	if p.Offset != 0 {
		q.Set("offset", strconv.FormatInt(p.Offset, 10))
	}
	if p.Timeout > 0 {
		q.Set("timeout", strconv.Itoa(p.Timeout))
	}
	if len(p.AllowedUpdates) > 0 {
		b, err := json.Marshal(p.AllowedUpdates)
		if err != nil {
			return nil, fmt.Errorf("getUpdates: marshal allowed_updates: %w", err)
		}
		q.Set("allowed_updates", string(b)) // e.g. ["message"]
	}

	data, err := c.doRequest(ctx, "getUpdates", q)
	if err != nil {
		return nil, fmt.Errorf("getUpdates request: %w", err)
	}

	res, perr := parseAPIResponse[[]Update](data)
	if perr != nil {
		return nil, fmt.Errorf("getUpdates parse: %w", perr)
	}
	return res, nil
}

// SendMessage sends a plain text message.
func (c *Client) SendMessage(ctx context.Context, p SendMessageParams) (*Message, error) {
	q := url.Values{}
	q.Set("chat_id", strconv.FormatInt(p.ChatID, 10))
	q.Set("text", p.Text)

	if p.ParseMode != "" {
		q.Set("parse_mode", p.ParseMode)
	}
	if p.DisableNotification {
		q.Set("disable_notification", "true")
	}
	if p.ReplyToMessageID != 0 {
		q.Set("reply_to_message_id", strconv.FormatInt(p.ReplyToMessageID, 10))
	}

	data, err := c.doRequest(ctx, "sendMessage", q)
	if err != nil {
		return nil, fmt.Errorf("sendMessage request: %w", err)
	}

	msg, perr := parseAPIResponse[*Message](data)
	if perr != nil {
		return nil, fmt.Errorf("sendMessage parse: %w", perr)
	}
	return msg, nil
}
