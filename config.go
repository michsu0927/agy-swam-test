package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds the runtime environment configuration for the service.
type Config struct {
	AppName string
	Env     string
	Port    int
}

// LoadConfig loads and validates configuration from environment variables.
func LoadConfig() (*Config, error) {
	appName := strings.TrimSpace(os.Getenv("APP_NAME"))
	if appName == "" {
		appName = "swarm-service"
	}

	appEnv := strings.TrimSpace(os.Getenv("APP_ENV"))
	if appEnv == "" {
		appEnv = "development"
	}
	switch appEnv {
	case "development", "staging", "production":
	default:
		return nil, fmt.Errorf("invalid APP_ENV %q: must be one of development, staging, production", appEnv)
	}

	port := 8080
	rawPort := strings.TrimSpace(os.Getenv("APP_PORT"))
	if rawPort != "" {
		parsedPort, err := strconv.Atoi(rawPort)
		if err != nil {
			return nil, fmt.Errorf("invalid APP_PORT %q: must be an integer: %w", rawPort, err)
		}
		if parsedPort < 1 || parsedPort > 65535 {
			return nil, fmt.Errorf("invalid APP_PORT %d: must be between 1 and 65535", parsedPort)
		}
		port = parsedPort
	}

	return &Config{
		AppName: appName,
		Env:     appEnv,
		Port:    port,
	}, nil
}

// Addr returns the HTTP listen address in the form ":<Port>".
func (c *Config) Addr() string {
	if c == nil || c.Port <= 0 {
		return ":8080"
	}
	return fmt.Sprintf(":%d", c.Port)
}
