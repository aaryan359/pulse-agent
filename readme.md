# Monitoring Agent

A lightweight, production-ready monitoring agent for tracking system and Docker container metrics.

## 🎯 What This Does

Collects system metrics and Docker container stats from your servers and sends them to your backend API. Built for **mobile-first monitoring** - check your servers from your phone at 2 AM.

## ✨ Features

- ✅ **System monitoring** - CPU, memory, disk usage
- ✅ **Docker monitoring** - Container stats, status, resource usage
- ✅ **Lightweight** - < 20MB Docker image
- ✅ **Simple setup** - One command installation
- ✅ **Auto-retry** - Handles network failures gracefully
- ✅ **Production-ready** - Used in real deployments

## 🚀 Quick Start

### For Users (Installing the Agent)

**One-line install:**
```bash
curl -fsSL https://yourapp.com/install.sh | bash
```

**Manual install:**
```bash
docker run -d \
  --name monitoring-agent \
  --restart unless-stopped \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -e AGENT_API_KEY="your_api_key" \
  -e AGENT_BACKEND_URL="https://api.yourapp.com" \
  yourregistry/monitoring-agent:latest
```

See [INSTALL.md](INSTALL.md) for detailed instructions.

### For Developers (Building the Agent)

**Prerequisites:**
- Go 1.21+
- Docker

**Build and run locally:**
```bash
# Clone repo
git clone <your-repo>
cd monitoring-agent

# Install dependencies
go mod download

# Run locally (requires Docker running)
export AGENT_API_KEY="test-key"
export AGENT_BACKEND_URL="http://localhost:8000"
go run cmd/agent/main.go

# Build Docker image
docker build -t monitoring-agent .

# Run in Docker
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -e AGENT_API_KEY="test-key" \
  monitoring-agent
```

## 📁 Project Structure

```
agent/
├── cmd/
│   └── agent/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration loading
│   ├── docker/
│   │   ├── client.go           # Docker client
│   │   └── stats.go            # Container stats
│   ├── system/
│   │   └── system.go           # System metrics
│   ├── collector/
│   │   └── collector.go        # Main collector
│   ├── sender/
│   │   └── sender.go           # Backend communication
│   ├── scheduler/
│   │   └── scheduler.go        # Collection scheduler
│   └── models/
│       └── payload.go          # Data structures
├── pkg/
│   └── logger/
│       └── logger.go           # Logging
├── Dockerfile
├── go.mod
└── README.md
```

## 🔧 Configuration

All configuration via environment variables:

```bash
AGENT_API_KEY          # Required - Your API key
AGENT_BACKEND_URL      # Backend API URL (default: https://api.yourapp.com)
AGENT_INTERVAL         # Collection interval in seconds (default: 10)
AGENT_SERVER_ID        # Server identifier (default: hostname)
AGENT_ENV              # Environment tag (default: production)
LOG_LEVEL              # info/debug/warn/error (default: info)
```

## 📊 Data Collected

### System Metrics
- CPU usage (total %)
- Memory usage (MB, %)
- Disk usage (GB, %)
- System uptime
- Host information

### Container Metrics (when Docker available)
- Container ID, name, image
- Status (running/stopped/exited)
- CPU usage per container
- Memory usage and limits
- Network I/O (RX/TX)

### Payload Example
```json
{
  "server_id": "production-web-01",
  "environment": "production",
  "timestamp": "2024-12-17T10:30:00Z",
  "system": {
    "hostname": "web-01",
    "cpu_percent": 25.3,
    "memory_used_mb": 2048,
    "memory_total_mb": 4096,
    "disk_used_gb": 50,
    "disk_total_gb": 100
  },
  "containers": [
    {
      "id": "abc123def456",
      "name": "nginx",
      "image": "nginx:latest",
      "state": "running",
      "cpu_percent": 5.2,
      "memory_usage_mb": 128
    }
  ],
  "container_count": 5
}
```

## 🏗️ Architecture

```
┌──────────────────────────┐
│        AGENT             │
│                          │
│  ┌───────────────┐       │
│  │  Collector    │────┐  │
│  └───────────────┘    │  │
│                        ▼  │
│  ┌───────────────┐  Metrics
│  │   Scheduler   │─────────▶ Backend
│  └───────────────┘          (HTTP)
│                             │
└──────────────────────────┘
```

**Design principles:**
- Simple loops, no magic
- Clear boundaries between components
- Easy to extend without rewrites
- Testable components
- Production-ready from day one

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with verbose output
go test -v ./...

# Test specific package
go test ./internal/collector

# With coverage
go test -cover ./...
```

## 📦 Deployment

### Build Docker image
```bash
docker build -t yourregistry/monitoring-agent:latest .
docker push yourregistry/monitoring-agent:latest
```

### Deploy install script
```bash
# Host install.sh on your domain
# Update AGENT_IMAGE in install.sh
# Users can then run:
curl -fsSL https://yourapp.com/install.sh | bash
```

### Update agents
```bash
# Push new image
docker build -t yourregistry/monitoring-agent:v1.1.0 .
docker tag yourregistry/monitoring-agent:v1.1.0 yourregistry/monitoring-agent:latest
docker push yourregistry/monitoring-agent:latest

# Users can update with:
docker pull yourregistry/monitoring-agent:latest
docker restart monitoring-agent
```

## 🔐 Security Considerations

- Agent requires **read-only** access to Docker socket
- Uses API key authentication (Bearer token)
- Communicates over HTTPS only
- No arbitrary command execution
- Limited to whitelisted operations
- No sensitive data stored locally

## 🚧 Roadmap

- [ ] Kubernetes support
- [ ] GPU metrics
- [ ] Custom metrics via plugins
- [ ] Compression for large payloads
- [ ] Local caching for offline periods
- [ ] Alert thresholds (client-side)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

## 📄 License

MIT License - see LICENSE file

## 🆘 Support

- Documentation: https://docs.yourapp.com
- Issues: https://github.com/yourorg/agent/issues
- Email: support@yourapp.com

---

Built with ❤️ for developers who want to monitor servers from their phones.
