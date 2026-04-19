package metrics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	c := New()
	c.RecordScan(5)
	c.RecordAlert()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	Handler(c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("unexpected content-type: %s", ct)
	}

	var out jsonSnapshot
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if out.ScansTotal != 1 {
		t.Errorf("expected 1 scan, got %d", out.ScansTotal)
	}
	if out.AlertsTotal != 1 {
		t.Errorf("expected 1 alert, got %d", out.AlertsTotal)
	}
	if out.UptimeSeconds < 0 {
		t.Error("uptime should be non-negative")
	}
}

func TestHandler_EmptyCollector(t *testing.T) {
	c := New()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	Handler(c).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
