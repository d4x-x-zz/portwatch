package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RotateManagerOptions configures the rotation manager.
type RotateManagerOptions struct {
	MaxSizeBytes int64
	MaxAge       time.Duration
	KeepLast     int
}

// DefaultRotateManagerOptions returns sensible defaults.
func DefaultRotateManagerOptions() RotateManagerOptions {
	return RotateManagerOptions{
		MaxSizeBytes: 10 * 1024 * 1024,
		MaxAge:       7 * 24 * time.Hour,
		KeepLast:     5,
	}
}

// MaybeRotate checks whether logPath needs rotation and, if so, renames it
// with a timestamp suffix and prunes old rotated files.
func MaybeRotate(logPath string, opts RotateManagerOptions) error {
	needs, err := NeedsRotation(logPath, DefaultRotateOptions{
		MaxSizeBytes: opts.MaxSizeBytes,
		MaxAge:       opts.MaxAge,
	})
	if err != nil || !needs {
		return err
	}

	timestamp := time.Now().UTC().Format("20060102T150405Z")
	dest := fmt.Sprintf("%s.%s", logPath, timestamp)
	if err := os.Rename(logPath, dest); err != nil {
		return fmt.Errorf("audit rotate: rename %s -> %s: %w", logPath, dest, err)
	}

	return pruneRotated(logPath, opts.KeepLast)
}

// pruneRotated removes old rotated files, keeping only the most recent keepLast.
func pruneRotated(logPath string, keepLast int) error {
	pattern := logPath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("audit rotate: glob %s: %w", pattern, err)
	}

	sort.Strings(matches) // lexicographic order == chronological for our timestamp format

	if len(matches) <= keepLast {
		return nil
	}

	for _, old := range matches[:len(matches)-keepLast] {
		if err := os.Remove(old); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("audit rotate: remove %s: %w", old, err)
		}
	}
	return nil
}
