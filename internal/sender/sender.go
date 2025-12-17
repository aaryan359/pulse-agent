// internal/sender/sender.go
package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	"pulse_agent/internal/config"
	"pulse_agent/internal/models"
)

type Sender struct {
	cfg    *config.Config
	client *http.Client
}

func New(cfg *config.Config) *Sender {
	return &Sender{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *Sender) Send(ctx context.Context, payload *models.Payload) error {
	jsonData, err := json.Marshal(payload)

	print("data to be send", jsonData)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// endpoint := fmt.Sprintf("%s/api/metrics", s.cfg.BackendURL)
	// req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	return fmt.Errorf("failed to create request: %w", err)
	// }

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.cfg.APIKey))
	// req.Header.Set("User-Agent", "monitoring-agent/1.0")

	// Retry logic with exponential backoff
	// maxRetries := 3
	// for attempt := 0; attempt < maxRetries; attempt++ {
	// 	resp, err := s.client.Do(req)
	// 	if err != nil {
	// 		logger.Warn("Send attempt %d failed: %v", attempt+1, err)
	// 		if attempt < maxRetries-1 {
	// 			time.Sleep(time.Duration(attempt+1) * 2 * time.Second)
	// 			continue
	// 		}
	// 		return fmt.Errorf("all retry attempts failed: %w", err)
	// 	}
	// 	defer resp.Body.Close()

	// 	body, _ := io.ReadAll(resp.Body)

	// 	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
	// 		logger.Debug("Metrics sent successfully")
	// 		return nil
	// 	}

	// 	if resp.StatusCode == 401 || resp.StatusCode == 403 {
	// 		return fmt.Errorf("authentication failed: %s", string(body))
	// 	}

	// 	if resp.StatusCode >= 500 && attempt < maxRetries-1 {
	// 		logger.Warn("Server error (attempt %d): %d - %s", attempt+1, resp.StatusCode, string(body))
	// 		time.Sleep(time.Duration(attempt+1) * 2 * time.Second)
	// 		continue
	// 	}

	// 	return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	// }

	return nil
}
