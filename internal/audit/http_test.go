package audit_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/monitor"
)

func TestHandler_Empty(t *testing.T) {
	l := audit.NewLog(tempLog(t))
	rec := httptest.NewRecorder()
	audit.Handler(l)(rec, httptest.NewRequest(http.MethodGet, "/audit", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []audit.Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty array")
	}
}

func TestHandler_WithEntries(t *testing.T) {
	l := audit.NewLog(tempLog(t))
	_ = l.Record(monitor.Diff{Opened: map[uint16]struct{}{443: {}}})
	rec := httptest.NewRecorder()
	audit.Handler(l)(rec, httptest.NewRequest(http.MethodGet, "/audit", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("unexpected content-type: %s", ct)
	}
	var entries []audit.Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}
