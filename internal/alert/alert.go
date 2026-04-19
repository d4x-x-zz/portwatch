// Package alert provides alerting mechanisms for port changes.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Alerter sends notifications when port changes are detected.
type Alerter interface {
	Notify(diff monitor.Diff) error
}

// LogAlerter writes alerts to a writer (default: stderr).
type LogAlerter struct {
	Out    io.Writer
	Prefix string
}

// NewLogAlerter creates a LogAlerter writing to stderr.
func NewLogAlerter() *LogAlerter {
	return &LogAlerter{
		Out:    os.Stderr,
		Prefix: "[portwatch]",
	}
}

// Notify writes a formatted alert message for the given diff.
func (a *LogAlerter) Notify(diff monitor.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}
	timestamp := time.Now().Format(time.RFC3339)
	_, err := fmt.Fprintf(a.Out, "%s %s port change detected at %s\n",
		a.Prefix, timestamp, diff.String())
	return err
}
