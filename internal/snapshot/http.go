package snapshot

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Handler returns an http.Handler that exposes snapshot listing and diff endpoints.
func Handler(dir string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/snapshots", listHandler(dir))
	mux.HandleFunc("/snapshots/diff", diffHandler(dir))
	return mux
}

func listHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := List(dir)
		if err != nil {
			http.Error(w, "failed to list snapshots: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	}
}

func diffHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n := 2
		if v := r.URL.Query().Get("n"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed >= 2 {
				n = parsed
			}
		}
		entries, err := List(dir)
		if err != nil || len(entries) < n {
			http.Error(w, "not enough snapshots", http.StatusNotFound)
			return
		}
		old, err := Load(entries[len(entries)-n].Path)
		if err != nil {
			http.Error(w, "failed to load old snapshot", http.StatusInternalServerError)
			return
		}
		new_, err := Load(entries[len(entries)-1].Path)
		if err != nil {
			http.Error(w, "failed to load new snapshot", http.StatusInternalServerError)
			return
		}
		d := Diff(old, new_)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(d)
	}
}
