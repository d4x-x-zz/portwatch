package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// WebhookAlerter sends port change notifications to an HTTP endpoint.
type WebhookAlerter struct {
	url    string
	client *http.Client
}

type webhookPayload struct {
	Timestamp string   `json:"timestamp"`
	Opened    []string `json:"opened"`
	Closed    []string `json:"closed"`
}

// NewWebhookAlerter creates a WebhookAlerter that posts to the given URL.
// If client is nil, a default client with a 5s timeout is used.
func NewWebhookAlerter(url string, client *http.Client) *WebhookAlerter {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &WebhookAlerter{url: url, client: client}
}

// Notify sends a JSON payload to the configured webhook URL if there are changes.
func (w *WebhookAlerter) Notify(diff monitor.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}

	opened := toStringSlice(diff.Opened)
	closed := toStringSlice(diff.Closed)

	payload := webhookPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Opened:    opened,
		Closed:    closed,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func toStringSlice(ports []int) []string {
	s := make([]string, len(ports))
	for i, p := range ports {
		s[i] = fmt.Sprintf("%d", p)
	}
	return s
}
