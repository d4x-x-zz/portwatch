package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func TestLogAlerter_Notify_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	a := &alert.LogAlerter{Out: &buf, Prefix: "[test]"}

	diff := monitor.Diff{
		Opened: []int{8080},
		Closed: []int{},
	}

	if err := a.Notify(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[test]") {
		t.Errorf("expected prefix in output, got: %s", out)
	}
	if !strings.Contains(out, "port change detected") {
		t.Errorf("expected change message in output, got: %s", out)
	}
}

func TestLogAlerter_Notify_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	a := &alert.LogAlerter{Out: &buf, Prefix: "[test]"}

	diff := monitor.Diff{
		Opened: []int{},
		Closed: []int{},
	}

	if err := a.Notify(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}

func TestNewLogAlerter_Defaults(t *testing.T) {
	a := alert.NewLogAlerter()
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
	if a.Prefix != "[portwatch]" {
		t.Errorf("unexpected prefix: %s", a.Prefix)
	}
}
