// Package notify provides alerting middleware for portwatch.
//
// Backoff strategies control how long the retry alerter waits between
// successive attempts. Two strategies are available:
//
//   - "exponential": delay doubles each attempt, capped at a maximum.
//   - "constant": every attempt uses the same fixed delay.
//
// Use DefaultExponentialBackoff for sensible production defaults.
package notify
