package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name:    "default configuration when no env vars set",
			envVars: map[string]string{},
			want: &Config{
				AppName: "swarm-service",
				Env:     "development",
				Port:    8080,
			},
			wantErr: false,
		},
		{
			name: "custom valid development configuration",
			envVars: map[string]string{
				"APP_NAME": "custom-app",
				"APP_ENV":  "development",
				"APP_PORT": "3000",
			},
			want: &Config{
				AppName: "custom-app",
				Env:     "development",
				Port:    3000,
			},
			wantErr: false,
		},
		{
			name: "custom valid staging configuration",
			envVars: map[string]string{
				"APP_NAME": "staging-app",
				"APP_ENV":  "staging",
				"APP_PORT": "8443",
			},
			want: &Config{
				AppName: "staging-app",
				Env:     "staging",
				Port:    8443,
			},
			wantErr: false,
		},
		{
			name: "custom valid production configuration",
			envVars: map[string]string{
				"APP_NAME": "prod-app",
				"APP_ENV":  "production",
				"APP_PORT": "80",
			},
			want: &Config{
				AppName: "prod-app",
				Env:     "production",
				Port:    80,
			},
			wantErr: false,
		},
		{
			name: "whitespace-only APP_NAME falls back to default",
			envVars: map[string]string{
				"APP_NAME": "   ",
				"APP_ENV":  "development",
				"APP_PORT": "8080",
			},
			want: &Config{
				AppName: "swarm-service",
				Env:     "development",
				Port:    8080,
			},
			wantErr: false,
		},
		{
			name: "trimmed environment variable values",
			envVars: map[string]string{
				"APP_NAME": "  trimmed-service  ",
				"APP_ENV":  "  staging  ",
				"APP_PORT": " 9090 ",
			},
			want: &Config{
				AppName: "trimmed-service",
				Env:     "staging",
				Port:    9090,
			},
			wantErr: false,
		},
		{
			name: "boundary port 1 is valid",
			envVars: map[string]string{
				"APP_PORT": "1",
			},
			want: &Config{
				AppName: "swarm-service",
				Env:     "development",
				Port:    1,
			},
			wantErr: false,
		},
		{
			name: "boundary port 65535 is valid",
			envVars: map[string]string{
				"APP_PORT": "65535",
			},
			want: &Config{
				AppName: "swarm-service",
				Env:     "development",
				Port:    65535,
			},
			wantErr: false,
		},
		{
			name: "invalid APP_ENV returns error",
			envVars: map[string]string{
				"APP_ENV": "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non-numeric APP_PORT returns error",
			envVars: map[string]string{
				"APP_PORT": "abc",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "out-of-range port 0 returns error",
			envVars: map[string]string{
				"APP_PORT": "0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative port -5 returns error",
			envVars: map[string]string{
				"APP_PORT": "-5",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "out-of-range port 65536 returns error",
			envVars: map[string]string{
				"APP_PORT": "65536",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("APP_NAME", "")
			t.Setenv("APP_ENV", "")
			t.Setenv("APP_PORT", "")

			for k, v := range tc.envVars {
				t.Setenv(k, v)
			}

			got, err := LoadConfig()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("LoadConfig() expected error, got nil (config: %+v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("LoadConfig() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("LoadConfig() returned nil config without error")
			}
			if got.AppName != tc.want.AppName || got.Env != tc.want.Env || got.Port != tc.want.Port {
				t.Errorf("LoadConfig() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestConfigAddr(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want string
	}{
		{
			name: "valid port 8080",
			cfg:  &Config{AppName: "swarm-service", Env: "development", Port: 8080},
			want: ":8080",
		},
		{
			name: "custom port 3000",
			cfg:  &Config{AppName: "swarm-service", Env: "production", Port: 3000},
			want: ":3000",
		},
		{
			name: "nil config falls back to :8080",
			cfg:  nil,
			want: ":8080",
		},
		{
			name: "zero port falls back to :8080",
			cfg:  &Config{AppName: "swarm-service", Env: "development", Port: 0},
			want: ":8080",
		},
		{
			name: "negative port falls back to :8080",
			cfg:  &Config{AppName: "swarm-service", Env: "development", Port: -100},
			want: ":8080",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.cfg.Addr(); got != tc.want {
				t.Errorf("Config.Addr() = %q, want %q", got, tc.want)
			}
		})
	}
}
