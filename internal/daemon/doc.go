// Package daemon implements the main run loop for portwatch.
//
// It periodically scans the configured port range, compares the results
// against the previously persisted state, fires alerts on any changes, and
// then saves the new state to disk.
//
// # Architecture
//
// The daemon coordinates three main components:
//
//   - A scanner that probes TCP/UDP ports in the configured range.
//   - A history store that persists port state between runs.
//   - An alerter that notifies on newly opened or closed ports.
//
// The scan interval and port range are controlled via the [config.Config]
// struct. If the context passed to [Daemon.Run] is cancelled, the daemon
// completes any in-progress scan before returning.
//
// # Typical usage
//
//	cfg, _ := config.Load("portwatch.toml")
//	store, _ := history.NewStore("/var/lib/portwatch/state.json")
//	alerter := alert.NewLogAlerter(nil)
//	d := daemon.New(cfg, store, alerter)
//	d.Run(ctx)
package daemon
