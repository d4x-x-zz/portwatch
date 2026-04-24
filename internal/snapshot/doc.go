// Package snapshot provides functionality for saving, loading, listing,
// pruning, and comparing port scan snapshots.
//
// Snapshots are point-in-time records of which ports were observed open
// during a scan. They are stored as JSON files in a configurable directory
// and can be used to audit historical port state or detect drift between
// two points in time.
//
// Use Diff to compare two snapshot entries and obtain a structured
// description of which ports were opened or closed between them.
package snapshot
