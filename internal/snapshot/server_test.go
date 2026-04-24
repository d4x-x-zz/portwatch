package snapshot_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func freeSnapshotAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func TestSnapshotServer_DisabledDoesNotListen(t *testing.T) {
	dir := tempDir(t)
	cfg := snapshot.ServerConfig{Addr: freeSnapshotAddr(t), Dir: dir, Enabled: false}
	s := snapshot.NewServer(cfg)
	if err := s.Start(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := s.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}
}

func TestSnapshotServer_StartsAndResponds(t *testing.T) {
	dir := tempDir(t)
	addr := freeSnapshotAddr(t)
	cfg := snapshot.ServerConfig{Addr: addr, Dir: dir, Enabled: true}
	s := snapshot.NewServer(cfg)
	if err := s.Start(); err != nil {
		t.Fatalf("start error: %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	resp, err := http.Get(fmt.Sprintf("http://%s/snapshots", addr))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if err := s.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}
}
