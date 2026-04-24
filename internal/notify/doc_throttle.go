// Package notify provides alerting middleware for portwatch.
//
// ThrottledAlerter wraps an Alerter and ensures notifications are not sent
// more frequently than a configured window allows. Excess notifications are
// buffered and flushed when the window expires or Flush is called explicitly.
package notify
