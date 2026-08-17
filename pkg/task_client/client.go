package task_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const apiPrefix = "/api/v1"

type Config struct {
	Host string `default:"localhost" envconfig:"HTTP_CLIENT_HOST"`
	Port string `default:"8080"      envconfig:"HTTP_CLIENT_PORT"`
}

type Client struct {
	client http.Client
	host   string
	token  string
}

func New(c Config) *Client {
	return &Client{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		host: net.JoinHostPort(c.Host, c.Port),
	}
}

func (c *Client) WithToken(token string) *Client {
	clone := *c
	clone.token = token

	return &clone
}

func (c *Client) Token() string {
	return c.token
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, in, out any) error {
	var body io.Reader

	if in != nil {
		raw, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshalling request body: %w", err)
		}
		body = bytes.NewReader(raw)
	}

	u := url.URL{Scheme: "http", Host: c.host, Path: path}
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return &APIError{StatusCode: resp.StatusCode, Message: errorMessage(raw)}
	}

	if out == nil || len(raw) == 0 {
		return nil
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("unmarshalling response (%s %s): %w", method, path, err)
	}

	return nil
}

func errorMessage(raw []byte) string {
	var body struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(raw, &body); err == nil && body.Error != "" {
		return body.Error
	}

	var plain string
	if err := json.Unmarshal(raw, &plain); err == nil && plain != "" {
		return plain
	}

	return string(raw)
}

func (c *Client) Ready(ctx context.Context) error {
	return c.do(ctx, http.MethodGet, "/ready", nil, nil, nil)
}

func (c *Client) WaitReady(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error

	for {
		lastErr = c.Ready(ctx)
		if lastErr == nil {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("service is not ready after %s: %w", timeout, lastErr)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
