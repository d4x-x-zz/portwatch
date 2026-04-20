package health

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordScan_IncrementsCount(t *testing.T) {
	st := NewState()
	st.RecordScan(nil)
	st.RecordScan(nil)
	if st.ScanCount != 2 {
		t.Fatalf("expected 2, got %d", st.ScanCount)
	}
}

func TestRecordScan_StoresError(t *testing.T) {
	st := NewState()
	st.RecordScan(errors.New("boom"))
	if st.LastError != "boom" {
		t.Fatalf("unexpected error string: %q", st.LastError)
	}
}

func TestRecordScan_ClearsError(t *testing.T) {
	st := NewState()
	st.RecordScan(errors.New("boom"))
	st.RecordScan(nil)
	if st.LastError != "" {
		t.Fatalf("expected empty error, got %q", st.LastError)
	}
}

func TestStatusHandler_JSON(t *testing.T) {
	st := NewState()
	st.RecordScan(nil)

	rec := httptest.NewRecorder()
	StatusHandler(st).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/status", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if _, ok := body["scan_count"]; !ok {
		t.Fatal("missing scan_count field")
	}
}

func TestStatusHandler_ContentType(t *testing.T) {
	st := NewState()
	rec := httptest.NewRecorder()
	StatusHandler(st).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/status", nil))
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
}
