package config

// AlertConfig holds configuration for all alerting backends.
type AlertConfig struct {
	Log     LogAlertConfig     `toml:"log"`
	Webhook WebhookAlertConfig `toml:"webhook"`
	Slack   SlackAlertConfig   `toml:"slack"`
}

// LogAlertConfig configures the log-based alerter.
type LogAlertConfig struct {
	Enabled bool `toml:"enabled"`
}

// WebhookAlertConfig configures the HTTP webhook alerter.
type WebhookAlertConfig struct {
	Enabled bool   `toml:"enabled"`
	URL     string `toml:"url"`
}

// SlackAlertConfig configures the Slack alerter.
type SlackAlertConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	Username string `toml:"username"`
}

// AnyEnabled reports whether at least one alerter is enabled.
func (a AlertConfig) AnyEnabled() bool {
	return a.Log.Enabled || a.Webhook.Enabled || a.Slack.Enabled
}

// validateAlerts checks that enabled alerters have required fields set.
func validateAlerts(a AlertConfig) error {
	if a.Webhook.Enabled && a.Webhook.URL == "" {
		return &ValidationError{Field: "alert.webhook.url", Msg: "url is required when webhook alerter is enabled"}
	}
	if a.Slack.Enabled && a.Slack.URL == "" {
		return &ValidationError{Field: "alert.slack.url", Msg: "url is required when slack alerter is enabled"}
	}
	return nil
}

// ValidationError describes a configuration validation failure.
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Msg
}
