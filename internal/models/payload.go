// internal/models/payload.go
package models

import "time"

type Payload struct {
	ServerID       string            `json:"server_id"`
	Environment    string            `json:"environment"`
	Timestamp      time.Time         `json:"timestamp"`
	System         *SystemMetric     `json:"system"`
	Containers     []ContainerMetric `json:"containers"`
	ContainerCount int               `json:"container_count"`
}

type SystemMetric struct {
	Hostname      string  `json:"hostname"`
	OS            string  `json:"os"`
	Platform      string  `json:"platform"`
	Arch          string  `json:"arch"`
	Uptime        uint64  `json:"uptime"`
	CPUCores      int     `json:"cpu_cores"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryTotalMB int     `json:"memory_total_mb"`
	MemoryUsedMB  int     `json:"memory_used_mb"`
	MemoryPercent float64 `json:"memory_percent"`
	DiskTotalGB   int     `json:"disk_total_gb"`
	DiskUsedGB    int     `json:"disk_used_gb"`
	DiskPercent   float64 `json:"disk_percent"`
}

type ContainerMetric struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	State         string    `json:"state"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	CPUPercent    float64   `json:"cpu_percent"`
	MemoryUsageMB int       `json:"memory_usage_mb"`
	MemoryLimitMB int       `json:"memory_limit_mb"`
	NetworkRxMB   float64   `json:"network_rx_mb"`
	NetworkTxMB   float64   `json:"network_tx_mb"`
}
