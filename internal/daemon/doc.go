// Package daemon implements the main run loop for portwatch.
//
// It periodically scans the configured port range, compares the results
// against the previously persisted state, fires alerts on any changes, and
// then saves the new state to disk.
//
// Typical usage:
//
//	cfg, _ := config.Load("portwatch.toml")
//	store, _ := history.NewStore("/var/lib/portwatch/state.json")
//	alerter := alert.NewLogAlerter(nil)
//	d := daemon.New(cfg, store, alerter)
//	d.Run(ctx)
package daemon
