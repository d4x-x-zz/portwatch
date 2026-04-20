// Package notify provides alerting middleware for portwatch.
//
// RateLimiter wraps an Alerter and suppresses duplicate notifications
// within a configurable cooldown window. Two diffs are considered equal
// when their opened and closed port sets are identical.
//
// Use NewRateLimiter to construct one; it is safe for concurrent use.
package notify
