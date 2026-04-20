// Circuit breaker for the notify pipeline.
//
// The circuit breaker wraps an Alerter and tracks consecutive failures.
// Once the failure threshold is reached the circuit opens and all calls
// are rejected immediately without forwarding to the downstream alerter.
//
// After a configurable cooldown the circuit moves to half-open: the next
// call is attempted. Success closes the circuit; failure reopens it.
package notify
