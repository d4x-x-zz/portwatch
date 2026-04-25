// Package notify provides alerting middleware for portwatch.
//
// # Debouncer
//
// Debouncer delays forwarding a diff until no new diffs have arrived
// within a configurable quiet period. This prevents alert storms when
// many ports change state in rapid succession.
//
// Use NewDebouncer to wrap any Alerter. Call Flush to immediately
// forward any pending diff without waiting for the quiet period.
//
// # Example
//
//	base := notify.NewLogAlerter(logger)
//	d := notify.NewDebouncer(base, 5*time.Second)
//	defer d.Flush()
//
// Diffs arriving within 5 seconds of each other are coalesced into a
// single alert. Once the quiet period elapses the merged diff is
// forwarded to the wrapped Alerter.
package notify
