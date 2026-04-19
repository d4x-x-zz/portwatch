package metrics

import (
	"encoding/json"
	"net/http"
	"time"
)

type jsonSnapshot struct {
	ScansTotal    int64     `json:"scans_total"`
	AlertsTotal   int64     `json:"alerts_total"`
	LastScanAt    time.Time `json:"last_scan_at,omitempty"`
	LastScanPorts int       `json:"last_scan_ports"`
	UpSince       time.Time `json:"up_since"`
	UptimeSeconds float64   `json:"uptime_seconds"`
}

// Handler returns an http.Handler that serves current metrics as JSON.
func Handler(c *Collector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := c.Snapshot()
		out := jsonSnapshot{
			ScansTotal:    s.ScansTotal,
			AlertsTotal:   s.AlertsTotal,
			LastScanAt:    s.LastScanAt,
			LastScanPorts: s.LastScanPorts,
			UpSince:       s.UpSince,
			UptimeSeconds: time.Since(s.UpSince).Seconds(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	})
}
