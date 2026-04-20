// Package audit provides a persistent, append-only log of port change events.
// Each diff recorded by the daemon is written as a JSON line to a file on disk.
// The HTTP handler exposes the full history as a JSON array for inspection.
package audit
