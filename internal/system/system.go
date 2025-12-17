// internal/system/system.go
package system

import (
	"context"
	"runtime"

	"pulse_agent/internal/models"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type Collector struct{}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) GetSystemStats(ctx context.Context) (*models.SystemMetric, error) {
	metric := &models.SystemMetric{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// Host info
	hostInfo, err := host.InfoWithContext(ctx)
	if err == nil {
		metric.Hostname = hostInfo.Hostname
		metric.Platform = hostInfo.Platform
		metric.Uptime = hostInfo.Uptime
	}

	// CPU stats
	cpuPercent, err := cpu.PercentWithContext(ctx, 0, false)
	if err == nil && len(cpuPercent) > 0 {
		metric.CPUPercent = cpuPercent[0]
	}

	cpuCount, err := cpu.CountsWithContext(ctx, true)
	if err == nil {
		metric.CPUCores = cpuCount
	}

	// Memory stats
	memInfo, err := mem.VirtualMemoryWithContext(ctx)
	if err == nil {
		metric.MemoryTotalMB = int(memInfo.Total / 1024 / 1024)
		metric.MemoryUsedMB = int(memInfo.Used / 1024 / 1024)
		metric.MemoryPercent = memInfo.UsedPercent
	}

	// Disk stats
	diskInfo, err := disk.UsageWithContext(ctx, "/")
	if err == nil {
		metric.DiskTotalGB = int(diskInfo.Total / 1024 / 1024 / 1024)
		metric.DiskUsedGB = int(diskInfo.Used / 1024 / 1024 / 1024)
		metric.DiskPercent = diskInfo.UsedPercent
	}

	return metric, nil
}
