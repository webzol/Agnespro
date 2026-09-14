// Package agnes is a thin client for the Agnes AI OpenAI-compatible
// API gateway (text, image, video) plus the video polling endpoint.
package agnes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client wraps an *http.Client with a base URL and API key. It is safe
// for concurrent use.
type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		HTTP:    &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *Client) do(ctx context.Context, method, p string, body any, out any) (*http.Response, error) {
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+p, rd)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		defer resp.Body.Close()
		raw, _ := io.ReadAll(resp.Body)
		return resp, fmt.Errorf("agnes: %s %s -> %d: %s", method, p, resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decode %s %s: %w", method, p, err)
		}
	}
	return resp, nil
}

func (c *Client) doAbsolute(ctx context.Context, method, u string, out any) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		defer resp.Body.Close()
		raw, _ := io.ReadAll(resp.Body)
		return resp, fmt.Errorf("agnes: %s %s -> %d: %s", method, u, resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decode %s %s: %w", method, u, err)
		}
	}
	return resp, nil
}

func (c *Client) Ping(ctx context.Context) error {
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	_, err := c.do(ctx, http.MethodPost, "/chat/completions", map[string]any{
		"model":     "agnes-2.5-flash",
		"messages":  []map[string]string{{"role": "user", "content": "ping"}},
		"max_tokens": 8,
	}, &resp)
	return err
}
