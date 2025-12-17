// internal/docker/stats.go
package docker

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"pulse_agent/internal/models"
	"pulse_agent/pkg/logger"

	"github.com/docker/docker/api/types/container"
)

func (c *Client) GetContainerStats(ctx context.Context) ([]models.ContainerMetric, error) {
	if c.cli == nil {
		return []models.ContainerMetric{}, nil
	}

	containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	metrics := make([]models.ContainerMetric, 0, len(containers))

	for _, ctr := range containers {
		metric := models.ContainerMetric{
			ID:        ctr.ID[:12],
			Name:      ctr.Names[0][1:],
			Image:     ctr.Image,
			State:     ctr.State,
			Status:    ctr.Status,
			CreatedAt: time.Unix(ctr.Created, 0),
		}

		if ctr.State == "running" {
			stats, err := c.getSingleContainerStats(ctx, ctr.ID)
			if err != nil {
				logger.Warn("Failed to get stats for %s: %v", metric.Name, err)
				continue
			}

			metric.CPUPercent = stats.CPUPercent
			metric.MemoryUsageMB = stats.MemoryUsageMB
			metric.MemoryLimitMB = stats.MemoryLimitMB
			metric.NetworkRxMB = stats.NetworkRxMB
			metric.NetworkTxMB = stats.NetworkTxMB
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (c *Client) getSingleContainerStats(ctx context.Context, containerID string) (*models.ContainerMetric, error) {
	resp, err := c.cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stats container.StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil && err != io.EOF {
		return nil, err
	}

	metric := &models.ContainerMetric{}

	// CPU %
	cpuDelta := float64(
		stats.CPUStats.CPUUsage.TotalUsage -
			stats.PreCPUStats.CPUUsage.TotalUsage,
	)

	systemDelta := float64(
		stats.CPUStats.SystemUsage -
			stats.PreCPUStats.SystemUsage,
	)

	if systemDelta > 0 && cpuDelta > 0 {
		metric.CPUPercent =
			(cpuDelta / systemDelta) *
				float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	// Memory
	metric.MemoryUsageMB = int(stats.MemoryStats.Usage / 1024 / 1024)
	metric.MemoryLimitMB = int(stats.MemoryStats.Limit / 1024 / 1024)

	// Network
	var rxBytes, txBytes uint64
	for _, net := range stats.Networks {
		rxBytes += net.RxBytes
		txBytes += net.TxBytes
	}

	metric.NetworkRxMB = float64(rxBytes) / 1024 / 1024
	metric.NetworkTxMB = float64(txBytes) / 1024 / 1024

	return metric, nil
}
