package main

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the JSON payload returned by GET /health.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Env     string `json:"env"`
}

// NewMux constructs an HTTP ServeMux with the /health endpoint registered.
func NewMux(cfg *Config) *http.ServeMux {
	if cfg == nil {
		cfg = &Config{
			AppName: "swarm-service",
			Env:     "development",
			Port:    8080,
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error":"method not allowed"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HealthResponse{
			Status:  "ok",
			Service: cfg.AppName,
			Env:     cfg.Env,
		})
	})

	return mux
}
