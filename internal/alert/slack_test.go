package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func TestSlackAlerter_Notify_WithChanges(t *testing.T) {
	var received slackPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := NewSlackAlerter(ts.URL, "bot")
	diff := monitor.Diff{Opened: []int{8080}, Closed: []int{}}
	if err := a.Notify(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Username != "bot" {
		t.Errorf("expected username 'bot', got %q", received.Username)
	}
	if received.Text == "" {
		t.Error("expected non-empty text")
	}
}

func TestSlackAlerter_Notify_NoChanges(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	a := NewSlackAlerter(ts.URL, "")
	diff := monitor.Diff{}
	if err := a.Notify(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call when no changes")
	}
}

func TestSlackAlerter_Notify_BadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	a := NewSlackAlerter(ts.URL, "portwatch")
	diff := monitor.Diff{Opened: []int{9090}}
	if err := a.Notify(diff); err == nil {
		t.Error("expected error on bad status")
	}
}

func TestNewSlackAlerter_DefaultUsername(t *testing.T) {
	a := NewSlackAlerter("http://example.com", "")
	if a.username != "portwatch" {
		t.Errorf("expected default username 'portwatch', got %q", a.username)
	}
}
