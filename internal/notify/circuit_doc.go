// Package notify provides alerting middleware for portwatch.
//
// CircuitBreaker wraps any Alerter and prevents cascading failures by
// tracking consecutive errors. Once the failure threshold is reached the
// circuit opens and all Notify calls return immediately with an error.
// After the configured cooldown period the circuit moves to half-open and
// allows a single probe call through; on success it closes again.
package notify
