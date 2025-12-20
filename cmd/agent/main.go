package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pulse_agent/internal/agent"
	"pulse_agent/internal/config"
	"pulse_agent/internal/scheduler"
	"pulse_agent/internal/sender"
	"pulse_agent/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()
	logger.Info("Starting monitoring agent...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}

	logger.Info(
		"Agent configured - Backend: %s, Interval: %v, Env: %s",
		cfg.BackendURL,
		cfg.Interval,
		cfg.Environment,
	)

	ctx := context.Background()

	// Try to register or verify existing server
	serverID, err := registerOrLoadServer(ctx, cfg)

	if err != nil {
		logger.Fatal("Failed to initialize agent: %v", err)
	}

	logger.Info("Agent ready with server ID: %s", serverID)

	// Attach server ID to config
	cfg.ServerID = serverID

	// Start scheduler (metrics collection + sending)
	sched := scheduler.New(cfg)
	go sched.Start()

	// Graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, stopping agent...")

	sched.Stop()
	time.Sleep(1 * time.Second)

	logger.Info("Agent stopped gracefully")
}
func registerOrLoadServer(ctx context.Context, cfg *config.Config) (string, error) {
	identity, err := agent.LoadServerIdentity()
	if err != nil {
		logger.Info("No identity found, registering agent...")
		return performRegistration(ctx, cfg)
	}

	// 🔑 API key changed → old server ID is INVALID
	if identity.APIKeyHash != agent.HashAPIKey(cfg.APIKey) {
		logger.Warn("API key changed, clearing old server identity")
		agent.ClearServerIdentity()
		return performRegistration(ctx, cfg)
	}

	logger.Info("Server identity verified")
	return identity.ServerID, nil
}

func performRegistration(ctx context.Context, cfg *config.Config) (string, error) {
	serverID, err := sender.RegisterAgent(ctx, cfg)
	if err != nil {
		return "", err
	}

	if err := agent.SaveServerIdentity(serverID, cfg.APIKey); err != nil {
		logger.Warn("Failed to save server identity: %v", err)
	}

	logger.Info("Agent registered successfully (server_id=%s)", serverID)
	return serverID, nil
}
