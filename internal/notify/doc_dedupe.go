// Package notify provides alerting middleware for portwatch.
//
// DedupeAlerter suppresses repeated notifications for the same diff
// until the diff changes or a reset occurs. This prevents alert storms
// when a port remains in an unexpected state across multiple scan cycles.
package notify
