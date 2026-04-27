// Package notify provides alerting middleware for portwatch.
//
// WindowLimiter caps the number of notifications sent within a sliding
// time window, dropping excess alerts to protect downstream systems
// from notification storms.
package notify
