// Package alert implements notification backends for portwatch.
//
// An Alerter receives a monitor.Diff and is responsible for delivering
// a notification through some channel (log, email, webhook, etc.).
//
// Currently implemented alerters:
//
//   - LogAlerter: writes human-readable messages to any io.Writer.
//
// Example usage:
//
//	alerter := alert.NewLogAlerter()
//	if err := alerter.Notify(diff); err != nil {
//		log.Printf("alert failed: %v", err)
//	}
package alert
