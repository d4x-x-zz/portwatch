// Package notify provides alerting middleware for portwatch.
//
// Debouncer delays forwarding a diff until no new diffs have arrived
// within a configurable quiet period. This prevents alert storms when
// many ports change state in rapid succession.
//
// Use NewDebouncer to wrap any Alerter. Call Flush to immediately
// forward any pending diff without waiting for the quiet period.
package notify
