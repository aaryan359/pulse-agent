// internal/scheduler/scheduler.go
package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"pulse_agent/internal/collector"
	"pulse_agent/internal/config"
	"pulse_agent/internal/sender"
	"pulse_agent/pkg/logger"
)

type Scheduler struct {
	cfg       *config.Config
	collector *collector.Collector
	sender    *sender.Sender
	stopChan  chan struct{}
}

func New(cfg *config.Config) *Scheduler {
	return &Scheduler{
		cfg:       cfg,
		collector: collector.New(cfg),
		sender:    sender.New(cfg),
		stopChan:  make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	logger.Info("Scheduler started with interval: %v", s.cfg.Interval)
	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	// Run immediately on start
	s.runCollection()

	for {
		select {
		case <-ticker.C:
			s.runCollection()
		case <-s.stopChan:
			logger.Info("Scheduler stopped")
			s.collector.Close()
			return
		}
	}
}

func (s *Scheduler) runCollection() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startTime := time.Now()

	payload, err := s.collector.Collect(ctx)
	if err != nil {
		logger.Error("Collection failed: %v", err)
		return
	}

	data, _ := json.MarshalIndent(payload, "", "  ")
	logger.Debug("FULL PAYLOAD:\n%s", string(data))

	err = s.sender.Send(ctx, payload)
	if err != nil {
		logger.Error("Send failed: %v", err)
		return
	}

	elapsed := time.Since(startTime)
	logger.Info("Collection cycle completed in %v (containers: %d, cpu: %.1f%%, memory: %.1f%%)",
		elapsed,
		payload.ContainerCount,
		payload.System.CPUPercent,
		payload.System.MemoryPercent,
	)
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}
