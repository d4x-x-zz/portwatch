package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/monitor"
)

func tempLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.jsonl")
}

func TestRecord_EmptyDiff(t *testing.T) {
	l := audit.NewLog(tempLog(t))
	if err := l.Record(monitor.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := l.ReadAll()
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecord_WritesEntry(t *testing.T) {
	l := audit.NewLog(tempLog(t))
	d := monitor.Diff{
		Opened: map[uint16]struct{}{8080: {}},
		Closed: map[uint16]struct{}{22: {}},
	}
	if err := l.Record(d); err != nil {
		t.Fatalf("record: %v", err)
	}
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	l := audit.NewLog("/tmp/portwatch_no_such_file_xyz.jsonl")
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Fatalf("expected nil entries")
	}
}

func TestRecord_Appends(t *testing.T) {
	path := tempLog(t)
	l := audit.NewLog(path)
	d := monitor.Diff{Opened: map[uint16]struct{}{9000: {}}}
	_ = l.Record(d)
	_ = l.Record(d)
	entries, _ := l.ReadAll()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	_ = os.Remove(path)
}
