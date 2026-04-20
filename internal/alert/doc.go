// Package alert provides alerting backends for portwatch.
//
// Alerters implement the Alerter interface and are notified when port
// changes are detected. Available implementations include:
//
//   - LogAlerter: writes alerts to a standard logger
//   - WebhookAlerter: posts JSON payloads to an HTTP endpoint
//   - SlackAlerter: posts formatted messages to a Slack incoming webhook
//   - MultiAlerter: fans out notifications to multiple alerters
//
// Use NewMultiAlerter to combine several alerters into one.
//
// # Error Handling
//
// Alerter implementations should return an error if the alert could not be
// delivered. MultiAlerter collects all errors from its children and returns
// them as a combined error so that a failure in one backend does not prevent
// others from being notified.
package alert
