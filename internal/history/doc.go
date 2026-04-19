// Package history manages persistent snapshots of open port scan results.
//
// A [Store] serialises the most recent [Snapshot] — a timestamped list of
// open ports — to a JSON file on disk. On the next daemon startup the
// snapshot is reloaded so that the monitor can detect changes that occurred
// while portwatch was not running.
package history
