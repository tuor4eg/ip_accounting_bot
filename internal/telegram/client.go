package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type response[T any] struct {
	Ok          bool   `json:"ok"`
	Result      T      `json:"result"`
	Description string `json:"description,omitempty"`
}

type Client struct {
	token   string
	baseURL string
	http    *http.Client
}

func New(token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &Client{
		token:   token,
		baseURL: "https://api.telegram.org",
		http:    httpClient,
	}
}

func (c *Client) buildURL(method string, q url.Values) string {
	u := fmt.Sprintf("%s/bot%s/%s", c.baseURL, c.token, method)

	if len(q) > 0 {
		u += "?" + q.Encode()
	}

	return u
}

func (c *Client) doRequest(ctx context.Context, method string, q url.Values) (*http.Response, error) {
	url := c.buildURL(method, q)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.http.Do(req)
}

func decodeJSON[T any](r io.Reader) (response[T], error) {
	var out response[T]

	err := json.NewDecoder(r).Decode(&out)

	return out, err
}

func parseAPIResponse[T any](res *http.Response) (T, error) {
	var zero T
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return zero, &APIError{
			Status:      res.StatusCode,
			Description: http.StatusText(res.StatusCode),
		}
	}

	out, err := decodeJSON[T](res.Body)

	if err != nil {
		return zero, err
	}

	if !out.Ok {
		description := out.Description

		if description != "" {
			description = "not ok"
		}

		return zero, &APIError{Description: description}
	}

	return out.Result, nil
}

//--------------------------------PUBLIC METHODS--------------------------------

func (c *Client) GetMe(ctx context.Context) (*User, error) {
	res, err := c.doRequest(ctx, "getMe", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return parseAPIResponse[*User](res)
}
