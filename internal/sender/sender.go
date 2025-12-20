package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"pulse_agent/internal/config"
	"pulse_agent/internal/models"
)

type Sender struct {
	cfg    *config.Config
	client *http.Client
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

var (
	ErrServerNotRegistered = errors.New("server_not_registered")
	ErrAuthFailed          = errors.New("authentication_failed")
	ErrInvalidResponse     = errors.New("invalid_response")
)

func New(cfg *config.Config) *Sender {
	return &Sender{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VerifyServerID checks if the server ID is still valid with current API key
func VerifyServerID(ctx context.Context, cfg *config.Config, serverID string) error {
	// Create a minimal test payload
	testPayload := &models.Payload{
		ServerID:    serverID,
		Environment: cfg.Environment,
		Timestamp:   time.Now(),
	}

	sender := New(cfg)
	err := sender.Send(ctx, testPayload)

	if errors.Is(err, ErrServerNotRegistered) {
		return fmt.Errorf("server ID not registered or API key changed")
	}

	if errors.Is(err, ErrAuthFailed) {
		return fmt.Errorf("authentication failed - API key may be invalid")
	}

	// Any other error might be temporary (network, etc) - consider it verified
	// The scheduler will handle retries for actual metric sending
	return nil
}

func (s *Sender) Send(ctx context.Context, payload *models.Payload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload failed: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/v1/agent/storeMetric", s.cfg.BackendURL)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.cfg.APIKey)
	req.Header.Set("User-Agent", "pulse-agent/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// ✅ Success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// 🔁 Server UUID invalid → re-register needed
	if resp.StatusCode == http.StatusConflict {
		var errResp struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		}

		_ = json.Unmarshal(body, &errResp)

		msg := strings.ToLower(errResp.Message)
		code := strings.ToLower(errResp.Code)

		if msg == "server not registered" ||
			msg == "server_not_registered" ||
			code == "server_not_registered" {
			return ErrServerNotRegistered
		}

		return fmt.Errorf("conflict: %s", string(body))
	}

	// ❌ Auth failure
	if resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("%w: %s", ErrAuthFailed, string(body))
	}

	// ❌ Bad request
	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("bad request: %s", string(body))
	}

	// ❌ Server error
	if resp.StatusCode >= 500 {
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("unexpected response %d: %s", resp.StatusCode, body)
}
