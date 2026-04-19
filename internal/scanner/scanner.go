package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port with its protocol and state.
type Port struct {
	Number   int
	Protocol string // "tcp" or "udp"
}

// Scanner probes ports on a given host.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a Scanner for the given host.
func New(host string, timeout time.Duration) *Scanner {
	return &Scanner{Host: host, Timeout: timeout}
}

// Scan checks a range of ports and returns the ones that are open.
func (s *Scanner) Scan(start, end int) ([]Port, error) {
	if start < 1 || end > 65535 || start > end {
		return nil, fmt.Errorf("invalid port range: %d-%d", start, end)
	}

	var open []Port
	for port := start; port <= end; port++ {
		addr := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err == nil {
			conn.Close()
			open = append(open, Port{Number: port, Protocol: "tcp"})
		}
	}
	return open, nil
}
