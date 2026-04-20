// Package audit records a rolling log of port change events.
package audit

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Entry is a single audit record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Opened    []uint16  `json:"opened,omitempty"`
	Closed    []uint16  `json:"closed,omitempty"`
}

// Log is a thread-safe, file-backed audit log.
type Log struct {
	mu   sync.Mutex
	path string
}

// NewLog returns a Log that persists entries to path.
func NewLog(path string) *Log {
	return &Log{path: path}
}

// Record appends an entry for the given diff (no-op when diff is empty).
func (l *Log) Record(d monitor.Diff) error {
	if len(d.Opened) == 0 && len(d.Closed) == 0 {
		return nil
	}
	e := Entry{
		Timestamp: time.Now().UTC(),
		Opened:    toSlice(d.Opened),
		Closed:    toSlice(d.Closed),
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(e)
}

// ReadAll returns all entries stored in the log file.
func (l *Log) ReadAll() ([]Entry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.Open(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return entries, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func toSlice(m map[uint16]struct{}) []uint16 {
	if len(m) == 0 {
		return nil
	}
	s := make([]uint16, 0, len(m))
	for p := range m {
		s = append(s, p)
	}
	return s
}
