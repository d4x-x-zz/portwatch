// Package alert provides alerting mechanisms for portwatch.
//
// Alerters are notified when port changes are detected by the monitor.
// The package currently provides a log-based alerter that writes
// structured output to a standard logger.
//
// Usage:
//
//	alerter := alert.NewLogAlerter(nil) // nil uses default logger
//	alerter.Notify(diff)
//
// Custom alerters can be implemented by satisfying the Alerter interface.
package alert
