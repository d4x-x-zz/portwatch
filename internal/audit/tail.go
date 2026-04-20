package audit

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const defaultTailN = 20

// TailHandler returns the last N audit log entries via HTTP.
// The caller may specify ?n=<count> to override the default.
func TailHandler(log *Log) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n := defaultTailN
		if raw := r.URL.Query().Get("n"); raw != "" {
			if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
				n = parsed
			}
		}

		all, err := log.ReadAll()
		if err != nil {
			http.Error(w, "failed to read audit log", http.StatusInternalServerError)
			return
		}

		start := len(all) - n
		if start < 0 {
			start = 0
		}
		slice := all[start:]

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(slice); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
