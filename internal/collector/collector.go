// internal/collector/collector.go
package collector

import (
	"context"
	"time"

	"pulse_agent/internal/config"
	"pulse_agent/internal/docker"
	"pulse_agent/internal/models"
	"pulse_agent/internal/system"
	"pulse_agent/pkg/logger"
)

type Collector struct {
	cfg          *config.Config
	dockerClient *docker.Client
	systemClient *system.Collector
}

func New(cfg *config.Config) *Collector {
	dockerClient, err := docker.NewClient()
	if err != nil {
		logger.Warn("Docker not available, system metrics only")
	}

	return &Collector{
		cfg:          cfg,
		dockerClient: dockerClient,
		systemClient: system.NewCollector(),
	}
}

func (c *Collector) Collect(ctx context.Context) (*models.Payload, error) {
	payload := &models.Payload{
		ServerID:    c.cfg.ServerID,
		Environment: c.cfg.Environment,
		Timestamp:   time.Now(),
	}

	// Collect system stats
	systemStats, err := c.systemClient.GetSystemStats(ctx)
	if err != nil {
		logger.Error("Failed to collect system stats: %v", err)
	} else {
		payload.System = systemStats
	}

	// Collect Docker stats if available
	if c.dockerClient != nil && c.dockerClient.IsAvailable() {
		containers, err := c.dockerClient.GetContainerStats(ctx)
		if err != nil {
			logger.Error("Failed to collect container stats: %v", err)
		} else {
			payload.Containers = payload.Containers
			payload.ContainerCount = len(containers)
		}
	}

	return payload, nil
}

func (c *Collector) Close() {
	if c.dockerClient != nil {
		c.dockerClient.Close()
	}
}
