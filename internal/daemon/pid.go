package daemon

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// WritePID writes the current process ID to the given file path.
func WritePID(path string) error {
	pid := os.Getpid()
	data := fmt.Sprintf("%d\n", pid)
	return os.WriteFile(path, []byte(data), 0o644)
}

// RemovePID deletes the PID file at the given path.
// It returns nil if the file does not exist.
func RemovePID(path string) error {
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// ReadPID reads and returns the PID stored in the given file.
func ReadPID(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid pid file %q: %w", path, err)
	}
	return pid, nil
}
