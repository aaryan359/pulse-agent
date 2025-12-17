# 🚀 Quick Start - Get Running in 5 Minutes

This guide gets you from zero to running agent in **5 minutes**.

---

## For Users: Install the Agent

### 1. Get your API key
```
Sign up at https://yourapp.com
Go to Settings → API Keys
Copy your key (starts with sk_live_...)
```

### 2. Run one command
```bash
curl -fsSL https://yourapp.com/install.sh | bash
```

### 3. Paste your API key when prompted
```
Enter your API key: sk_live_abc123...
```

### 4. Done! 🎉
Your server appears in the mobile app within 30 seconds.

---

## For Developers: Build and Test Locally

### 1. Clone and setup
```bash
git clone <your-repo>
cd monitoring-agent
cp .env.example .env
```

### 2. Edit .env file
```bash
# .env
AGENT_API_KEY=test-key-123
AGENT_BACKEND_URL=http://localhost:8000
LOG_LEVEL=debug
```

### 3. Run with Docker Compose (easiest)
```bash
docker-compose up
```

This starts:
- The agent
- A mock backend (httpbin)
- Example containers (nginx, redis)

### 4. Watch it work
```bash
# In another terminal, watch logs
docker-compose logs -f agent

# You should see:
# [INFO] Starting monitoring agent...
# [INFO] Docker client connected successfully
# [INFO] Collection cycle completed (containers: 3, cpu: 15.2%, memory: 35.8%)
```

### 5. Test manually (without Docker Compose)
```bash
# Install dependencies
go mod download

# Run locally
AGENT_API_KEY=test-key \
AGENT_BACKEND_URL=http://localhost:8000 \
go run cmd/agent/main.go
```

---

## Testing the Full Flow

### 1. Start a simple backend
```bash
# Terminal 1: Simple HTTP server that logs requests
python3 -m http.server 8000
```

### 2. Run the agent
```bash
# Terminal 2: Run agent pointing to local backend
AGENT_API_KEY=test-key \
AGENT_BACKEND_URL=http://localhost:8000/api/metrics \
go run cmd/agent/main.go
```

### 3. See metrics being sent
```
# Terminal 1 output:
127.0.0.1 - - [17/Dec/2024 10:30:15] "POST /api/metrics HTTP/1.1" 200 -
```

---

## Build Your Own Backend (5-minute version)

```javascript
// server.js - Minimal backend to receive metrics
const express = require('express');
const app = express();

app.use(express.json());

app.post('/api/metrics', (req, res) => {
  const { server_id, system, containers } = req.body;
  
  console.log(`📊 Metrics from ${server_id}:`);
  console.log(`  CPU: ${system.cpu_percent.toFixed(1)}%`);
  console.log(`  Memory: ${system.memory_percent.toFixed(1)}%`);
  console.log(`  Containers: ${containers.length}`);
  
  res.json({ success: true });
});

app.listen(8000, () => {
  console.log('🚀 Backend listening on http://localhost:8000');
});
```

Run it:
```bash
npm install express
node server.js
```

---

## Common Commands

```bash
# Build Docker image
docker build -t monitoring-agent .

# Run Docker image
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -e AGENT_API_KEY=test-key \
  -e AGENT_BACKEND_URL=http://host.docker.internal:8000 \
  monitoring-agent

# View logs
docker logs -f monitoring-agent

# Stop agent
docker stop monitoring-agent

# Remove agent
docker rm monitoring-agent
```

---

## Makefile Commands (if you have make)

```bash
make help              # Show all commands
make build             # Build binary
make run               # Run locally (needs AGENT_API_KEY env var)
make test              # Run tests
make docker-build      # Build Docker image
make docker-run        # Run in Docker (needs AGENT_API_KEY)
make example-local     # Run with example config
```

---

## What Data Looks Like

This is what the agent sends to your backend every 10 seconds:

```json
{
  "server_id": "my-server",
  "environment": "production",
  "timestamp": "2024-12-17T10:30:00Z",
  "system": {
    "hostname": "web-01",
    "cpu_percent": 25.3,
    "memory_used_mb": 2048,
    "memory_total_mb": 4096,
    "memory_percent": 50.0,
    "disk_used_gb": 45,
    "disk_total_gb": 100,
    "disk_percent": 45.0
  },
  "containers": [
    {
      "id": "abc123",
      "name": "nginx",
      "image": "nginx:latest",
      "state": "running",
      "cpu_percent": 2.5,
      "memory_usage_mb": 128,
      "memory_limit_mb": 512
    }
  ],
  "container_count": 3
}
```

---

## Troubleshooting

### "Docker not available"
```bash
# Make sure Docker is running
docker ps

# Check Docker socket exists
ls -la /var/run/docker.sock
```

### "Connection refused"
```bash
# Check your backend is running
curl http://localhost:8000

# Check firewall isn't blocking
sudo ufw status
```

### "API key invalid"
```bash
# Make sure you're setting the env var
echo $AGENT_API_KEY

# Try running with explicit key
AGENT_API_KEY=your-key-here go run cmd/agent/main.go
```

---

## Next Steps

1. ✅ **Got it running?** Great! Now build your backend
2. 📱 **Build mobile app** to display the metrics
3. 🚀 **Deploy** using the DEPLOYMENT.md guide
4. 📊 **Add features** like alerts, graphs, history

---

## Need Help?

- 📖 Read the full [README.md](README.md)
- 🔧 Check [INSTALL.md](INSTALL.md) for detailed installation
- 🚀 See [DEPLOYMENT.md](DEPLOYMENT.md) for production setup
- 🐛 Open an issue on GitHub

---

**That's it! You're ready to build your monitoring app. 🎉**

Time to make something awesome! 💪
