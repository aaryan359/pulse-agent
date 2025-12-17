// internal/docker/client.go
package docker

import (
	"context"
	"pulse_agent/pkg/logger"

	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx := context.Background()
	_, err = cli.Ping(ctx)
	if err != nil {
		logger.Warn("Docker not available: %v", err)
		return nil, err
	}

	logger.Info("Docker client connected successfully")
	return &Client{cli: cli}, nil
}

func (c *Client) Close() error {
	if c.cli != nil {
		return c.cli.Close()
	}
	return nil
}

func (c *Client) IsAvailable() bool {
	return c.cli != nil
}
