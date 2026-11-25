package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		cfg            *Config
		method         string
		path           string
		wantStatus     int
		wantAllow      string
		wantContentTyp string
		wantResp       *HealthResponse
	}{
		{
			name:           "GET /health with custom config returns 200 and JSON payload",
			cfg:            &Config{AppName: "my-service", Env: "production", Port: 9000},
			method:         http.MethodGet,
			path:           "/health",
			wantStatus:     http.StatusOK,
			wantContentTyp: "application/json",
			wantResp: &HealthResponse{
				Status:  "ok",
				Service: "my-service",
				Env:     "production",
			},
		},
		{
			name:           "GET /health with nil config falls back to default config",
			cfg:            nil,
			method:         http.MethodGet,
			path:           "/health",
			wantStatus:     http.StatusOK,
			wantContentTyp: "application/json",
			wantResp: &HealthResponse{
				Status:  "ok",
				Service: "swarm-service",
				Env:     "development",
			},
		},
		{
			name:       "POST /health returns 405 Method Not Allowed with Allow GET",
			cfg:        &Config{AppName: "swarm-service", Env: "development", Port: 8080},
			method:     http.MethodPost,
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantAllow:  http.MethodGet,
		},
		{
			name:       "PUT /health returns 405 Method Not Allowed with Allow GET",
			cfg:        &Config{AppName: "swarm-service", Env: "development", Port: 8080},
			method:     http.MethodPut,
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantAllow:  http.MethodGet,
		},
		{
			name:       "DELETE /health returns 405 Method Not Allowed with Allow GET",
			cfg:        &Config{AppName: "swarm-service", Env: "development", Port: 8080},
			method:     http.MethodDelete,
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantAllow:  http.MethodGet,
		},
		{
			name:       "GET /unknown returns 404 Not Found",
			cfg:        &Config{AppName: "swarm-service", Env: "development", Port: 8080},
			method:     http.MethodGet,
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := NewMux(tc.cfg)
			if mux == nil {
				t.Fatalf("NewMux() returned nil")
			}

			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.wantStatus {
				t.Errorf("status code = %d, want %d", res.StatusCode, tc.wantStatus)
			}

			if tc.wantAllow != "" {
				if gotAllow := res.Header.Get("Allow"); gotAllow != tc.wantAllow {
					t.Errorf("Allow header = %q, want %q", gotAllow, tc.wantAllow)
				}
			}

			if tc.wantContentTyp != "" {
				gotCT := res.Header.Get("Content-Type")
				if !strings.Contains(strings.ToLower(gotCT), tc.wantContentTyp) {
					t.Errorf("Content-Type header = %q, want it to contain %q", gotCT, tc.wantContentTyp)
				}
			}

			if tc.wantResp != nil {
				var got HealthResponse
				if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response JSON: %v", err)
				}
				if got != *tc.wantResp {
					t.Errorf("decoded response = %+v, want %+v", got, *tc.wantResp)
				}
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		wantAddr string
		wantSvc  string
		wantEnv  string
	}{
		{
			name:     "NewServer with custom config sets fields and timeouts",
			cfg:      &Config{AppName: "api-svc", Env: "staging", Port: 9090},
			wantAddr: ":9090",
			wantSvc:  "api-svc",
			wantEnv:  "staging",
		},
		{
			name:     "NewServer with nil config falls back safely to defaults",
			cfg:      nil,
			wantAddr: ":8080",
			wantSvc:  "swarm-service",
			wantEnv:  "development",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := NewServer(tc.cfg)
			if srv == nil {
				t.Fatalf("NewServer() returned nil")
			}
			if srv.Addr != tc.wantAddr {
				t.Errorf("srv.Addr = %q, want %q", srv.Addr, tc.wantAddr)
			}
			if srv.Handler == nil {
				t.Fatalf("srv.Handler is nil")
			}
			if srv.ReadHeaderTimeout <= 0 {
				t.Errorf("srv.ReadHeaderTimeout = %v, want > 0", srv.ReadHeaderTimeout)
			}
			if srv.ReadTimeout <= 0 {
				t.Errorf("srv.ReadTimeout = %v, want > 0", srv.ReadTimeout)
			}
			if srv.WriteTimeout <= 0 {
				t.Errorf("srv.WriteTimeout = %v, want > 0", srv.WriteTimeout)
			}
			if srv.IdleTimeout <= 0 {
				t.Errorf("srv.IdleTimeout = %v, want > 0", srv.IdleTimeout)
			}

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			srv.Handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("GET /health via srv.Handler status = %d, want %d", rec.Code, http.StatusOK)
			}
			var got HealthResponse
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode health response from srv.Handler: %v", err)
			}
			if got.Status != "ok" || got.Service != tc.wantSvc || got.Env != tc.wantEnv {
				t.Errorf("srv.Handler /health = %+v, want status=ok service=%q env=%q", got, tc.wantSvc, tc.wantEnv)
			}
		})
	}
}
