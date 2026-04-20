// Package health exposes a lightweight HTTP handler that reports whether
// the portwatch daemon is running and how long it has been up.
//
// Mount the handler on any ServeMux:
//
//	mux.Handle("/health", health.Handler(startTime))
package health
