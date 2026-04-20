package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "audit.log")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp log: %v", err)
	}
	return p
}

func TestNeedsRotation_MissingFile(t *testing.T) {
	ok, err := NeedsRotation("/nonexistent/audit.log", DefaultRotateOptions(), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected false for missing file")
	}
}

func TestNeedsRotation_UnderLimit(t *testing.T) {
	p := writeTempLog(t, "small")
	opts := RotateOptions{MaxBytes: 1024, MaxAge: time.Hour}
	ok, err := NeedsRotation(p, opts, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected false for small file")
	}
}

func TestNeedsRotation_ExceedsSize(t *testing.T) {
	p := writeTempLog(t, strings.Repeat("x", 100))
	opts := RotateOptions{MaxBytes: 50}
	ok, err := NeedsRotation(p, opts, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true when file exceeds MaxBytes")
	}
}

func TestNeedsRotation_ExceedsAge(t *testing.T) {
	p := writeTempLog(t, "data")
	opts := RotateOptions{MaxAge: time.Millisecond}
	time.Sleep(5 * time.Millisecond)
	ok, err := NeedsRotation(p, opts, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true when file exceeds MaxAge")
	}
}

func TestRotate_RenamesFile(t *testing.T) {
	p := writeTempLog(t, "entry\n")
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	archive, err := Rotate(p, now)
	if err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("original file should no longer exist")
	}
	if !strings.Contains(archive, "20240601T120000Z") {
		t.Fatalf("archive name %q missing timestamp", archive)
	}
	if _, err := os.Stat(archive); err != nil {
		t.Fatalf("archive not found: %v", err)
	}
}
