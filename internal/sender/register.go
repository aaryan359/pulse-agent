package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"pulse_agent/internal/config"
	"pulse_agent/pkg/logger"
)

type registerResponse struct {
	Data struct {
		ServerUUID string `json:"server_uuid"`
	} `json:"data"`
	Message string `json:"message"`
}

// RegisterAgent registers the server and returns server_uuid
func RegisterAgent(ctx context.Context, cfg *config.Config) (string, error) {
	payload := map[string]string{
		"hostname":    cfg.Hostname,
		"environment": cfg.Environment,
		"os":          cfg.OS,
		"arch":        cfg.Arch,
	}

	// TODO: separate logical server identity from agent UUID
	// to preserve server history across API key rotation

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal registration payload failed: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/v1/agent/register", cfg.BackendURL)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", fmt.Errorf("create registration request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", cfg.APIKey)
	req.Header.Set("User-Agent", "pulse-agent/1.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("registration request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	logger.Debug("Agent registration response status=%d body=%s",
		resp.StatusCode, string(responseBody))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"agent registration failed (%d): %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var result registerResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", fmt.Errorf("parse registration response failed: %w", err)
	}

	if result.Data.ServerUUID == "" {
		return "", fmt.Errorf("server_uuid not found in response")
	}

	return result.Data.ServerUUID, nil
}
