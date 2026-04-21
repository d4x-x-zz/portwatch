// Package notify provides alerting middleware for portwatch.
//
// RetryAlerter wraps an Alerter and retries failed Notify calls
// using a configurable backoff strategy. It is useful when the
// downstream alerting target (e.g. a webhook or Slack) may be
// temporarily unavailable.
//
// Use NewRetryAlerter to construct one, then compose it into a
// Pipeline alongside other middleware such as CircuitBreaker or
// ThrottledAlerter.
package notify
