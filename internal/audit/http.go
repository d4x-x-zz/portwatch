package audit

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.HandlerFunc that serves the audit log as JSON.
func Handler(l *Log) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := l.ReadAll()
		if err != nil {
			http.Error(w, "failed to read audit log", http.StatusInternalServerError)
			return
		}
		if entries == nil {
			entries = []Entry{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries) //nolint:errcheck
	}
}
