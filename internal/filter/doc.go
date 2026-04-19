// Package filter provides allow/deny rules for port alerting.
//
// Rules are evaluated in order; the first matching rule wins.
// If no rules match, the port is allowed by default.
//
// Rule format:
//
//	"allow:80"        - allow single port
//	"deny:8000-9000"  - deny port range
package filter
