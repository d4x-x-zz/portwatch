// Package notify provides alerting middleware for portwatch.
//
// DedupeAlerter
//
// DedupeAlerter wraps any Alerter and drops consecutive notifications that
// carry an identical diff. This is useful when a scan cycle produces the
// same change set multiple times in a row (e.g. a port that flaps and
// settles) and you want exactly one alert per unique event rather than a
// stream of repeats.
//
// It differs from RateLimiter in that there is no time window involved:
// suppression is purely content-based. A call to Reset() clears the
// memory so the next notification is always forwarded.
package notify
