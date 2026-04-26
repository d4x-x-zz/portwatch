// Package notify provides alerting middleware for portwatch.
//
// WindowLimiter wraps an Alerter and enforces a maximum number of
// notifications within a sliding time window. Once the limit is
// reached, subsequent calls are silently dropped until the window
// resets. This protects downstream systems from alert storms.
package notify
