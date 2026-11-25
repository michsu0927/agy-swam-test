package main

import (
	"log"
	"net/http"
	"time"
)

// NewServer constructs an HTTP server configured with timeouts and the service mux.
func NewServer(cfg *Config) *http.Server {
	if cfg == nil {
		cfg = &Config{
			AppName: "swarm-service",
			Env:     "development",
			Port:    8080,
		}
	}

	return &http.Server{
		Addr:              cfg.Addr(),
		Handler:           NewMux(cfg),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	srv := NewServer(cfg)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server exited with error: %v", err)
	}
}
