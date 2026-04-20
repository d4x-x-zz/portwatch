// Package monitor compares port scan results and detects changes.
package monitor

import (
	"fmt"
	"sort"
	"strings"
)

// PortSet is a set of open ports.
type PortSet map[int]struct{}

// Diff holds the result of comparing two PortSets.
type Diff struct {
	Opened []int
	Closed []int
}

// HasChanges returns true if any ports opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// String returns a human-readable summary of the diff.
func (d Diff) String() string {
	var parts []string
	if len(d.Opened) > 0 {
		parts = append(parts, fmt.Sprintf("opened: %v", d.Opened))
	}
	if len(d.Closed) > 0 {
		parts = append(parts, fmt.Sprintf("closed: %v", d.Closed))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}

// NewPortSet converts a slice of port numbers into a PortSet.
func NewPortSet(ports []int) PortSet {
	ps := make(PortSet, len(ports))
	for _, p := range ports {
		ps[p] = struct{}{}
	}
	return ps
}

// Ports returns a sorted slice of all ports in the set.
func (ps PortSet) Ports() []int {
	ports := make([]int, 0, len(ps))
	for p := range ps {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	return ports
}

// Compare returns a Diff between a previous and current PortSet.
// The Opened and Closed slices in the result are sorted in ascending order.
func Compare(prev, curr PortSet) Diff {
	var d Diff
	for p := range curr {
		if _, ok := prev[p]; !ok {
			d.Opened = append(d.Opened, p)
		}
	}
	for p := range prev {
		if _, ok := curr[p]; !ok {
			d.Closed = append(d.Closed, p)
		}
	}
	sort.Ints(d.Opened)
	sort.Ints(d.Closed)
	return d
}
