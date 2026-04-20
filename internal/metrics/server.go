package metrics

import (
	"context"
	"log"
	"net/http"
	"time"
)

// ServerConfig holds the parameters needed to start the metrics HTTP server.
type ServerConfig struct {
	Addr string
	Path string
}

// Server wraps an http.Server exposing the metrics handler at a configured path.
type Server struct {
	cfg    ServerConfig
	inner  *http.Server
	collector Collector
}

// NewServer creates a Server using the given config and collector.
func NewServer(cfg ServerConfig, c Collector) *Server {
	mux := http.NewServeMux()
	mux.Handle(cfg.Path, Handler(c))
	return &Server{
		cfg: cfg,
		collector: c,
		inner: &http.Server{
			Addr:    cfg.Addr,
			Handler: mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

// Start begins listening in a background goroutine.
func (s *Server) Start() {
	go func() {
		log.Printf("metrics: listening on %s%s", s.cfg.Addr, s.cfg.Path)
		if err := s.inner.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("metrics: server error: %v", err)
		}
	}()
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	return s.inner.Shutdown(ctx)
}
