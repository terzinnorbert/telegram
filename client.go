package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// TelegramClient communicates with the Telegram Bot API.
type TelegramClient struct {
	Token      string
	BaseURL    string
	HTTPClient *http.Client
}

// NewTelegramClient creates a new Telegram API client.
func NewTelegramClient(token string) *TelegramClient {
	return &TelegramClient{
		Token:   token,
		BaseURL: "https://api.telegram.org/bot",
		HTTPClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

func (c *TelegramClient) endpoint(method string) string {
	return fmt.Sprintf("%s%s/%s", c.BaseURL, c.Token, method)
}

func executeWithRetry[T any](ctx context.Context, c *TelegramClient, buildReq func() (*http.Request, error)) (*APIResponse[T], error) {
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := buildReq()
		if err != nil {
			return nil, err
		}
		req = req.WithContext(ctx)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
				continue
			}
			return nil, fmt.Errorf("network request failed: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			var apiResp APIResponse[T]
			_ = json.Unmarshal(body, &apiResp)
			waitTime := 1 * time.Second
			if apiResp.Parameters != nil && apiResp.Parameters.RetryAfter > 0 {
				waitTime = time.Duration(apiResp.Parameters.RetryAfter) * time.Second
			}
			if attempt < maxRetries-1 {
				select {
				case <-time.After(waitTime):
					continue
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
		}

		var apiResp APIResponse[T]
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to parse Telegram response (HTTP %d): %s", resp.StatusCode, string(body))
		}

		if !apiResp.OK {
			return nil, fmt.Errorf("Telegram API error (%d): %s", apiResp.ErrorCode, apiResp.Description)
		}

		return &apiResp, nil
	}

	return nil, fmt.Errorf("exceeded max retries")
}

// GetMe retrieves basic bot account information.
func (c *TelegramClient) GetMe(ctx context.Context) (*User, error) {
	reqFunc := func() (*http.Request, error) {
		return http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint("getMe"), nil)
	}

	resp, err := executeWithRetry[User](ctx, c, reqFunc)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

// SendMessage sends a text message to the specified chat ID.
func (c *TelegramClient) SendMessage(ctx context.Context, chatID string, text string, parseMode string) (*Message, error) {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	if parseMode != "" && parseMode != "plain" {
		payload["parse_mode"] = parseMode
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	reqFunc := func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint("sendMessage"), bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}

	resp, err := executeWithRetry[Message](ctx, c, reqFunc)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

// SendDocument uploads a document/file to the specified chat ID.
func (c *TelegramClient) SendDocument(ctx context.Context, chatID string, filename string, fileReader io.Reader, caption string, parseMode string) (*Message, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("chat_id", chatID); err != nil {
		return nil, fmt.Errorf("failed to write chat_id field: %w", err)
	}

	if caption != "" {
		if err := writer.WriteField("caption", caption); err != nil {
			return nil, fmt.Errorf("failed to write caption field: %w", err)
		}
	}

	if parseMode != "" && parseMode != "plain" {
		if err := writer.WriteField("parse_mode", parseMode); err != nil {
			return nil, fmt.Errorf("failed to write parse_mode field: %w", err)
		}
	}

	part, err := writer.CreateFormFile("document", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create document multipart field: %w", err)
	}

	if _, err := io.Copy(part, fileReader); err != nil {
		return nil, fmt.Errorf("failed to copy file payload: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize multipart writer: %w", err)
	}

	contentType := writer.FormDataContentType()

	reqFunc := func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint("sendDocument"), bytes.NewReader(body.Bytes()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)
		return req, nil
	}

	resp, err := executeWithRetry[Message](ctx, c, reqFunc)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

// GetUpdates fetches incoming updates.
func (c *TelegramClient) GetUpdates(ctx context.Context, offset int64, limit int) ([]Update, error) {
	reqFunc := func() (*http.Request, error) {
		url := c.endpoint("getUpdates")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		q := req.URL.Query()
		if limit > 0 {
			q.Set("limit", strconv.Itoa(limit))
		}
		if offset != 0 {
			q.Set("offset", strconv.FormatInt(offset, 10))
		}
		req.URL.RawQuery = q.Encode()
		return req, nil
	}

	resp, err := executeWithRetry[[]Update](ctx, c, reqFunc)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}
