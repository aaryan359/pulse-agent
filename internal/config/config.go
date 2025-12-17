// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIKey      string
	BackendURL  string
	Interval    time.Duration
	ServerID    string
	Environment string
}

func Load() (*Config, error) {
	env := getEnv("AGENT_ENV", "production")
	devMode := env == "development"

	apiKey := getEnv("AGENT_API_KEY", "")
	if apiKey == "" && !devMode {
		return nil, fmt.Errorf("AGENT_API_KEY is required")
	}

	backendURL := getEnv("AGENT_BACKEND_URL", "https://api.yourapp.com")
	intervalSec := getEnvInt("AGENT_INTERVAL", 1)
	serverID := getEnv("AGENT_SERVER_ID", generateServerID())
	environment := getEnv("AGENT_ENV", "production")

	return &Config{
		APIKey:      apiKey,
		BackendURL:  backendURL,
		Interval:    time.Duration(intervalSec) * time.Second,
		ServerID:    serverID,
		Environment: environment,
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func generateServerID() string {
	// Try to get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown-server"
	}
	return hostname
}
