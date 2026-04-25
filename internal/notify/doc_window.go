// Package notify provides alerting middleware for portwatch.
//
// Window limiter
//
// NewWindowLimiter wraps an Alerter and enforces a maximum number of
// notifications within a sliding time window. Once the cap is reached
// further calls are silently dropped until the window resets.
//
//	base := alert.NewLogAlerter(log.Default())
//	limited := notify.NewWindowLimiter(base, notify.WindowOptions{
//	    Window:   time.Minute,
//	    MaxCalls: 5,
//	})
package notify
