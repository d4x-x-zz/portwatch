// Package notify provides middleware-style wrappers around alert.Alerter
// implementations.
//
// Currently includes:
//
//   - RateLimiter: suppresses repeated identical notifications within a
//     configurable cooldown window, reducing noise when ports flap.
package notify
