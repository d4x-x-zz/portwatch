// Package notify provides alerting middleware for portwatch.
//
// DedupeAlerter suppresses consecutive identical diffs so that the
// downstream alerter is only called when the port state actually
// changes between successive scans.
package notify
