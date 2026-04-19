package scanner

import (
	"net"
	"testing"
	"time"
)

// startListener opens a TCP listener on a random port and returns it.
func startListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	return ln, ln.Addr().(*net.TCPAddr).Port
}

func TestScan_DetectsOpenPort(t *testing.T) {
	ln, port := startListener(t)
	defer ln.Close()

	s := New("127.0.0.1", 500*time.Millisecond)
	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 || ports[0].Number != port {
		t.Errorf("expected port %d to be open, got %v", port, ports)
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := New("127.0.0.1", 500*time.Millisecond)
	_, err := s.Scan(100, 10)
	if err == nil {
		t.Error("expected error for invalid range, got nil")
	}
}

func TestScan_ClosedPort(t *testing.T) {
	// Bind then immediately close to get a port number unlikely to be reused.
	ln, port := startListener(t)
	ln.Close()

	s := New("127.0.0.1", 200*time.Millisecond)
	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected no open ports, got %v", ports)
	}
}
