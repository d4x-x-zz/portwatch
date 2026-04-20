package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/health"
)

func TestHandler_ReturnsOK(t *testing.T) {
	upSince := time.Now().Add(-5 * time.Minute)
	h := health.Handler(upSince)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var s health.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !s.OK {
		t.Error("expected ok=true")
	}
	if s.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
}

func TestHandler_ContentType(t *testing.T) {
	h := health.Handler(time.Now())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("unexpected content-type: %s", ct)
	}
}
