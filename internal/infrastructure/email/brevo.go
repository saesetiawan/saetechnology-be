package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-platform-core/internal/config"
	"go-platform-core/internal/domain/email"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.brevo.com/v3/smtp/email"

type BrevoClient struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewBrevoClient(cfg *config.Config) email.EmailSender {
	return &BrevoClient{
		APIKey: cfg.BrevoApiKey,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *BrevoClient) SendEmail(payload email.EmailRequest) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	bodyString := string(bodyBytes)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("failed request: %s", bodyString)
	}
	return nil
}
