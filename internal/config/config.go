package config

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey      string
	BackendURL  string
	Interval    time.Duration
	ServerID    string
	Hostname    string
	Environment string
	OS          string
	Arch        string
}

func Load() (*Config, error) {
	// Load .env if present (no error in prod)
	_ = godotenv.Load()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}

	env := getEnv("AGENT_ENV", "production")
	isDev := env == "development"

	cfg := &Config{
		Environment: env,
		BackendURL:  getEnv("AGENT_BACKEND_URL", "http://localhost:3000"),
		Hostname:    hostname,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
	}

	// API Key (required in non-dev)
	cfg.APIKey = getEnv("AGENT_API_KEY", "")
	if cfg.APIKey == "" && !isDev {
		return nil, errors.New("AGENT_API_KEY is required in production")
	}

	// Interval (supports 1s, 5s, 1m)
	intervalRaw := getEnv("AGENT_INTERVAL", "10s")
	interval, err := parseInterval(intervalRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid AGENT_INTERVAL: %w", err)
	}
	cfg.Interval = interval

	// Validate backend URL
	if cfg.BackendURL == "" {
		return nil, errors.New("AGENT_BACKEND_URL is required")
	}

	return cfg, nil
}

/* -------------------- helpers -------------------- */

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseInterval(value string) (time.Duration, error) {
	// Try duration format first: "1s", "500ms", "1m"
	if d, err := time.ParseDuration(value); err == nil {
		if d < 1*time.Second {
			return 0, fmt.Errorf("interval must be at least 1 second")
		}
		return d, nil
	}

	// Fallback: plain seconds "5"
	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("must be duration (e.g. 5s, 1m) or seconds")
	}

	if seconds < 1 {
		return 0, fmt.Errorf("interval must be at least 1 second")
	}

	return time.Duration(seconds) * time.Second, nil
}
