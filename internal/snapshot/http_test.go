package snapshot_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func TestListHandler_Empty(t *testing.T) {
	dir := tempDir(t)
	h := snapshot.Handler(dir)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/snapshots", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []snapshot.Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty list, got %d", len(entries))
	}
}

func TestListHandler_WithEntries(t *testing.T) {
	dir := tempDir(t)
	if err := snapshot.Save(dir, []int{80, 443}, time.Now()); err != nil {
		t.Fatal(err)
	}
	h := snapshot.Handler(dir)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/snapshots", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []snapshot.Entry
	json.NewDecoder(rec.Body).Decode(&entries)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestDiffHandler_NotEnoughSnapshots(t *testing.T) {
	dir := tempDir(t)
	h := snapshot.Handler(dir)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/snapshots/diff", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestDiffHandler_ReturnsDiff(t *testing.T) {
	dir := tempDir(t)
	now := time.Now()
	snapshot.Save(dir, []int{80}, now)
	snapshot.Save(dir, []int{80, 443}, now.Add(time.Second))
	h := snapshot.Handler(dir)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/snapshots/diff", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("unexpected content-type: %s", ct)
	}
}
