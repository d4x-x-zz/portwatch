package snapshot

import (
	"context"
	"net/http"
	"time"
)

// ServerConfig holds configuration for the snapshot HTTP server.
type ServerConfig struct {
	Addr    string
	Dir     string
	Enabled bool
}

// Server serves snapshot endpoints over HTTP.
type Server struct {
	cfg ServerConfig
	srv *http.Server
}

// NewServer creates a new snapshot Server.
func NewServer(cfg ServerConfig) *Server {
	h := Handler(cfg.Dir)
	return &Server{
		cfg: cfg,
		srv: &http.Server{
			Addr:    cfg.Addr,
			Handler: h,
		},
	}
}

// Start begins listening. It returns immediately; call Shutdown to stop.
func (s *Server) Start() error {
	if !s.cfg.Enabled {
		return nil
	}
	go s.srv.ListenAndServe() //nolint:errcheck
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if !s.cfg.Enabled {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
