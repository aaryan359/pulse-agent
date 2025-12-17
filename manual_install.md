# Monitoring Agent - Installation Guide

## 🚀 Quick Install (Recommended)

### One-line install:
```bash
curl -fsSL https://yourapp.com/install.sh | bash
```

You'll be prompted for your API key during installation.

### With environment variables:
```bash
curl -fsSL https://yourapp.com/install.sh | \
  API_KEY=your_api_key_here \
  BACKEND_URL=https://api.yourapp.com \
  bash
```

---

## 📦 Manual Installation

### Prerequisites
- Docker installed and running
- API key from your dashboard

### Step 1: Pull the agent image
```bash
docker pull yourregistry/monitoring-agent:latest
```

### Step 2: Run the agent
```bash
docker run -d \
  --name monitoring-agent \
  --restart unless-stopped \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -e AGENT_API_KEY="your_api_key_here" \
  -e AGENT_BACKEND_URL="https://api.yourapp.com" \
  -e AGENT_SERVER_ID="$(hostname)" \
  -e AGENT_INTERVAL="10" \
  -e LOG_LEVEL="info" \
  yourregistry/monitoring-agent:latest
```

### Step 3: Verify it's running
```bash
docker ps | grep monitoring-agent
docker logs monitoring-agent
```

You should see log output like:
```
[2024-12-17 10:30:00] [INFO] Starting monitoring agent...
[2024-12-17 10:30:00] [INFO] Agent configured - Backend: https://api.yourapp.com
[2024-12-17 10:30:01] [INFO] Docker client connected successfully
[2024-12-17 10:30:01] [INFO] Scheduler started with interval: 10s
[2024-12-17 10:30:01] [INFO] Collection cycle completed (containers: 5, cpu: 25.3%, memory: 45.2%)
```

---

## ⚙️ Configuration Options

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `AGENT_API_KEY` | Yes | - | Your API key from dashboard |
| `AGENT_BACKEND_URL` | No | `https://api.yourapp.com` | Backend API URL |
| `AGENT_SERVER_ID` | No | hostname | Unique server identifier |
| `AGENT_INTERVAL` | No | `10` | Collection interval (seconds) |
| `AGENT_ENV` | No | `production` | Environment tag |
| `LOG_LEVEL` | No | `info` | Log level (debug/info/warn/error) |

### Example with all options:
```bash
docker run -d \
  --name monitoring-agent \
  --restart unless-stopped \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -e AGENT_API_KEY="sk_live_abc123..." \
  -e AGENT_BACKEND_URL="https://api.yourapp.com" \
  -e AGENT_SERVER_ID="production-web-01" \
  -e AGENT_INTERVAL="15" \
  -e AGENT_ENV="production" \
  -e LOG_LEVEL="info" \
  yourregistry/monitoring-agent:latest
```

---

## 🔧 Management Commands

### View logs
```bash
docker logs -f monitoring-agent
```

### Stop agent
```bash
docker stop monitoring-agent
```

### Start agent
```bash
docker start monitoring-agent
```

### Restart agent
```bash
docker restart monitoring-agent
```

### Update agent to latest version
```bash
docker stop monitoring-agent
docker rm monitoring-agent
docker pull yourregistry/monitoring-agent:latest
# Then run the installation command again
```

### Completely remove agent
```bash
docker stop monitoring-agent
docker rm monitoring-agent
docker rmi yourregistry/monitoring-agent:latest
```

---

## 🐛 Troubleshooting

### Agent not showing in app
1. Check agent logs: `docker logs monitoring-agent`
2. Verify API key is correct
3. Check backend URL is accessible
4. Ensure server has internet connection

### Docker socket permission error
```bash
# Add your user to docker group
sudo usermod -aG docker $USER
# Then logout and login again
```

### Agent crashes on start
```bash
# Check logs for error
docker logs monitoring-agent

# Common issues:
# - Invalid API key → Check your dashboard
# - Docker socket not accessible → Check -v mount
# - Backend unreachable → Check AGENT_BACKEND_URL
```

### High CPU usage
```bash
# Increase interval to reduce collection frequency
docker stop monitoring-agent
docker rm monitoring-agent
# Then run with AGENT_INTERVAL="30" (30 seconds)
```

---

## 📊 What Gets Monitored

### System Metrics
- CPU usage (per core and total)
- Memory usage and available
- Disk usage and capacity
- System uptime
- Host information

### Container Metrics (if Docker available)
- Container status (running/stopped)
- CPU usage per container
- Memory usage per container
- Network traffic (RX/TX)
- Container restart count

---

## 🔐 Security

- Agent runs with **read-only** Docker socket access
- Only collects metrics, **cannot modify** containers
- Uses secure HTTPS connection to backend
- API key authentication
- No sensitive data stored locally

---

## 🆘 Support

- Documentation: https://docs.yourapp.com
- Issues: https://github.com/yourorg/agent/issues
- Email: support@yourapp.com

---

## 📝 License

MIT License - See LICENSE file for details
