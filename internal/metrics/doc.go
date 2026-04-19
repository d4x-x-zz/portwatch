// Package metrics provides a lightweight, goroutine-safe collector for
// portwatch runtime statistics such as scan counts, alert counts, and
// uptime. Use Collector.Snapshot to retrieve a consistent point-in-time
// view without blocking the daemon's hot path.
package metrics
