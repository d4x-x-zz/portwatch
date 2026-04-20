package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions controls log rotation behaviour.
type RotateOptions struct {
	// MaxBytes is the maximum file size before rotation. Zero disables size-based rotation.
	MaxBytes int64
	// MaxAge is the maximum age of the log file before rotation. Zero disables age-based rotation.
	MaxAge time.Duration
}

// DefaultRotateOptions returns sensible rotation defaults.
func DefaultRotateOptions() RotateOptions {
	return RotateOptions{
		MaxBytes: 10 * 1024 * 1024, // 10 MiB
		MaxAge:   24 * time.Hour,
	}
}

// NeedsRotation reports whether the log file at path should be rotated
// given the provided options and the current time.
func NeedsRotation(path string, opts RotateOptions, now time.Time) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if opts.MaxBytes > 0 && info.Size() >= opts.MaxBytes {
		return true, nil
	}
	if opts.MaxAge > 0 && now.Sub(info.ModTime()) >= opts.MaxAge {
		return true, nil
	}
	return false, nil
}

// Rotate renames path to a timestamped archive name and returns the archive path.
// The caller is responsible for re-opening the log file after rotation.
func Rotate(path string, now time.Time) (string, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	stamp := now.UTC().Format("20060102T150405Z")
	archive := filepath.Join(dir, fmt.Sprintf("%s.%s", base, stamp))
	if err := os.Rename(path, archive); err != nil {
		return "", fmt.Errorf("audit rotate: %w", err)
	}
	return archive, nil
}
