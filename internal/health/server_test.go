package health

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func freeAddr(t *testing.T) string {
	t.Helper()
	// Use port 0 trick via a listener would be ideal, but for simplicity
	// pick an unlikely port for tests.
	return fmt.Sprintf("127.0.0.1:%d", 19876)
}

func TestServer_StartsAndResponds(t *testing.T) {
	srv := NewServer(freeAddr(t))
	if err := srv.Start(); err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	})

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + freeAddr(t) + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestServer_Shutdown(t *testing.T) {
	srv := NewServer("127.0.0.1:19877")
	if err := srv.Start(); err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error: %v", err)
	}
}
