// Package notify provides alerting middleware for portwatch.
//
// Pipeline chains multiple alerter wrappers together in a defined order,
// allowing operators to compose rate limiting, deduplication, throttling,
// circuit breaking, and retry behaviour declaratively from configuration.
package notify
