package oca

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-starter-app/config"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Client OCA client
type Client struct {
	httpClient   *http.Client
	maxRetries   int
	retryBackoff time.Duration
	email        *email
	whatsApp     *WhatsAppConfig
}

// ErrorResponse dari API
type ErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"errors"`
}

// NewClient Initialize http lient
func NewClient(cfg *config.Config) *Client {
	// set max retries (default 3)
	return &Client{
		httpClient:   &http.Client{Timeout: 15 * time.Second},
		maxRetries:   cfg.Oca().MaxRetries,
		retryBackoff: cfg.Oca().RetryBackoff,
		email: &email{
			baseUrl:     cfg.Oca().EmailBaseURL,
			authToken:   cfg.Oca().EmailAuthToken,
			senderEmail: cfg.Oca().EmailSenderEmail,
			senderName:  cfg.Oca().EmailSenderName,
			urlPrefix:   "/email/v2",
		},
		whatsApp: &WhatsAppConfig{
			BaseURL:   cfg.Oca().WhatsAppBaseURL,
			AuthToken: cfg.Oca().WhatsAppAuthToken,
		},
	}
}

// doRequest request handler with retry + rate limit handling
func (c *Client) doRequest(req *http.Request, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(c.retryBackoff * time.Duration(attempt+1))
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		err = resp.Body.Close()
		if err != nil {
			return err
		}

		// ✅ sukses
		if resp.StatusCode == http.StatusOK {
			if result != nil {
				if err := json.Unmarshal(body, result); err != nil {
					return err
				}
			}
			return nil
		}

		// ❌ 429 Too Many Requests → tunggu lalu retry
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := 5 * time.Second
			if val := resp.Header.Get("Retry-After"); val != "" {
				if sec, err := strconv.Atoi(val); err == nil {
					retryAfter = time.Duration(sec) * time.Second
				}
			}
			time.Sleep(retryAfter)
			continue
		}

		// ❌ 5xx → retry dengan backoff
		if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
			lastErr = fmt.Errorf("server error %d: %s", resp.StatusCode, string(body))
			time.Sleep(c.retryBackoff * time.Duration(attempt+1))
			continue
		}

		// ❌ error lain → langsung keluar
		var apiErr ErrorResponse
		if err := json.Unmarshal(body, &apiErr); err == nil && len(apiErr.Errors) > 0 {
			return errors.New(apiErr.Errors[0].Message)
		}
		return fmt.Errorf("unexpected error %d: %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("request failed after %d retries: %v", c.maxRetries, lastErr)
}
