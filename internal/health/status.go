package health

import (
	"encoding/json"
	"net/http"
	"time"
)

// StatusResponse is the JSON body returned by the status endpoint.
type StatusResponse struct {
	Status  string    `json:"status"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
}

// StatusHandler returns an http.HandlerFunc that reports daemon status.
func StatusHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := StatusResponse{
			Status:  "ok",
			Version: version,
			Time:    time.Now().UTC(),
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}
