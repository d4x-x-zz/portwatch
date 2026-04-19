// Package notify provides middleware wrappers around alert backends.
//
// Available middleware:
//   - RateLimiter  – suppresses duplicate diffs within a cooldown window.
//   - Debouncer    – waits for a quiet period before forwarding a diff.
//   - ThrottledAlerter – forwards at most one notification per interval,
//     queuing the latest pending diff for delivery via Flush.
package notify
