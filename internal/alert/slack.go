package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// SlackAlerter sends alerts to a Slack webhook URL.
type SlackAlerter struct {
	webhookURL string
	client     *http.Client
	username   string
}

type slackPayload struct {
	Username string `json:"username,omitempty"`
	Text     string `json:"text"`
}

// NewSlackAlerter creates a SlackAlerter that posts to the given Slack webhook URL.
func NewSlackAlerter(webhookURL, username string) *SlackAlerter {
	if username == "" {
		username = "portwatch"
	}
	return &SlackAlerter{
		webhookURL: webhookURL,
		username:   username,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends a Slack message if there are port changes.
func (s *SlackAlerter) Notify(diff monitor.Diff) error {
	if !diff.HasChanges() {
		return nil
	}

	payload := slackPayload{
		Username: s.username,
		Text:     fmt.Sprintf(":warning: Port change detected:\n%s", diff.String()),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
