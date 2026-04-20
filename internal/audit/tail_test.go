package audit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

func seedLog(t *testing.T, n int) *Log {
	t.Helper()
	path := tempLog(t)
	log, err := NewLog(path)
	if err != nil {
		t.Fatalf("NewLog: %v", err)
	}
	for i := 0; i < n; i++ {
		diff := monitor.Diff{
			Opened: []int{8000 + i},
		}
		if err := log.Record(fmt.Sprintf("seed-%d", i), diff, time.Now()); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	return log
}

func TestTailHandler_DefaultN(t *testing.T) {
	log := seedLog(t, 25)
	h := TailHandler(log)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit/tail", nil)
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) != 20 {
		t.Fatalf("expected 20 entries, got %d", len(entries))
	}
}

func TestTailHandler_CustomN(t *testing.T) {
	log := seedLog(t, 10)
	h := TailHandler(log)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit/tail?n=5", nil)
	h(rec, req)

	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) != 5 {
		t.Fatalf("expected 5, got %d", len(entries))
	}
}

func TestTailHandler_FewerThanN(t *testing.T) {
	log := seedLog(t, 3)
	h := TailHandler(log)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit/tail?n=20", nil)
	h(rec, req)

	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3, got %d", len(entries))
	}
}

func TestTailHandler_ContentType(t *testing.T) {
	log := seedLog(t, 1)
	h := TailHandler(log)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit/tail", nil)
	h(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}
