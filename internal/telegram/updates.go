package telegram

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// --- Public API methods ---

func (c *Client) GetUpdates(ctx context.Context, p GetUpdatesParams) ([]Update, error) {
	q := url.Values{}

	if p.Offset != 0 {
		q.Set("offset", strconv.FormatInt(p.Offset, 10))
	}

	if p.Timeout > 0 {
		q.Set("timeout", strconv.Itoa(p.Timeout))
	}

	if len(p.AllowedUpdates) > 0 {
		allowedUpdates, err := json.Marshal(p.AllowedUpdates)

		if err != nil {
			return nil, err
		}

		q.Set("allowed_updates", string(allowedUpdates))
	}

	data, err := c.doRequest(ctx, "getUpdates", q)

	if err != nil {
		return nil, err
	}

	return parseAPIResponse[[]Update](data)
}

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
		return nil, err
	}

	return parseAPIResponse[*Message](data)
}
