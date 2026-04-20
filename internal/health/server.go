package health

import (
	"context"
	"net/http"
	"time"
)

// Server wraps an HTTP server exposing the health endpoint.
type Server struct {
	addr string
	server *http.Server
}

// NewServer creates a health Server listening on addr.
func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/healthz", Handler())
	return &Server{
		addr: addr,
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

// Start begins serving in a background goroutine.
// It returns immediately; use Shutdown to stop.
func (s *Server) Start() error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return err
	case <-time.After(50 * time.Millisecond):
		return nil
	}
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
