package health

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// State holds runtime status reported by the health status endpoint.
type State struct {
	mu      sync.RWMutex
	Started time.Time
	LastScan time.Time
	LastError string
	ScanCount int64
}

// NewState returns a State initialised with the current time.
func NewState() *State {
	return &State{Started: time.Now()}
}

// RecordScan updates the last scan timestamp and increments the counter.
func (s *State) RecordScan(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastScan = time.Now()
	s.ScanCount++
	if err != nil {
		s.LastError = err.Error()
	} else {
		s.LastError = ""
	}
}

// StatusHandler returns an http.Handler that serialises State as JSON.
func StatusHandler(st *State) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st.mu.RLock()
		snap := struct {
			Started   time.Time `json:"started"`
			LastScan  time.Time `json:"last_scan"`
			LastError string    `json:"last_error,omitempty"`
			ScanCount int64     `json:"scan_count"`
		}{
			Started:   st.Started,
			LastScan:  st.LastScan,
			LastError: st.LastError,
			ScanCount: st.ScanCount,
		}
		st.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(snap)
	})
}
