// Package health provides a simple HTTP health-check endpoint.
package health

import (
	"encoding/json"
	"net/http"
	"time"
)

// Status holds the response payload for the health endpoint.
type Status struct {
	OK      bool      `json:"ok"`
	UpSince time.Time `json:"up_since"`
	Uptime  string    `json:"uptime"`
}

// Handler returns an http.Handler that reports service health.
func Handler(upSince time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := Status{
			OK:      true,
			UpSince: upSince,
			Uptime:  time.Since(upSince).Round(time.Second).String(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(status)
	})
}
