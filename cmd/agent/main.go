// cmd/agent/main.go
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"pulse_agent/internal/config"
	"pulse_agent/internal/scheduler"
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

	logger.Info("Agent configured - Backend: %s, Interval: %v", cfg.BackendURL, cfg.Interval)

	// Start scheduler
	sched := scheduler.New(cfg)
	go sched.Start()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, stopping agent...")
	sched.Stop()
	time.Sleep(1 * time.Second)
	logger.Info("Agent stopped gracefully")
}
