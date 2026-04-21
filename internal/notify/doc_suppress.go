// Package notify provides alerting middleware components.
//
// SuppressAlerter wraps an Alerter and suppresses repeated identical
// notifications after a configurable threshold. Once the suppression
// limit is reached, further identical diffs are silently dropped until
// either the diff changes (optionally) or Reset is called.
package notify
