package snapshot

import (
	"fmt"
	"sort"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// SnapshotDiff describes the changes between two snapshots.
type SnapshotDiff struct {
	OldTime time.Time
	NewTime time.Time
	Opened  []int
	Closed  []int
}

// Empty returns true when there are no changes between the two snapshots.
func (d SnapshotDiff) Empty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// String returns a human-readable summary of the diff.
func (d SnapshotDiff) String() string {
	if d.Empty() {
		return "no changes between snapshots"
	}
	return fmt.Sprintf("snapshot diff %s -> %s: +%d opened, -%d closed",
		d.OldTime.Format(time.RFC3339),
		d.NewTime.Format(time.RFC3339),
		len(d.Opened),
		len(d.Closed),
	)
}

// Diff compares two snapshot entries and returns the ports that were opened
// or closed between them.
func Diff(old, new Entry) SnapshotDiff {
	oldSet := monitor.NewPortSet(old.Ports)
	newSet := monitor.NewPortSet(new.Ports)

	md := monitor.Compare(oldSet, newSet)

	opened := sortedSlice(md.Opened)
	closed := sortedSlice(md.Closed)

	return SnapshotDiff{
		OldTime: old.Time,
		NewTime: new.Time,
		Opened:  opened,
		Closed:  closed,
	}
}

func sortedSlice(s map[int]struct{}) []int {
	out := make([]int, 0, len(s))
	for p := range s {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}
