# Complete Deployment Guide

## 🎯 Deployment Strategy

This guide covers deploying the agent from development to production.

---

## 📋 Pre-Deployment Checklist

- [ ] Backend API is running and accessible
- [ ] API authentication is configured
- [ ] Docker registry is set up (Docker Hub / GitHub Registry / Private)
- [ ] Install script is hosted and accessible
- [ ] Mobile app is ready to receive metrics
- [ ] Monitoring/alerting for the agent itself

---

## 1️⃣ Build and Push Agent Image

### Option A: Docker Hub
```bash
# Login to Docker Hub
docker login

# Build the image
docker build -t yourusername/monitoring-agent:latest .

# Tag with version
docker tag yourusername/monitoring-agent:latest yourusername/monitoring-agent:v1.0.0

# Push to registry
docker push yourusername/monitoring-agent:latest
docker push yourusername/monitoring-agent:v1.0.0
```

### Option B: GitHub Container Registry
```bash
# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Build the image
docker build -t ghcr.io/yourusername/monitoring-agent:latest .

# Push
docker push ghcr.io/yourusername/monitoring-agent:latest
```

### Option C: Private Registry
```bash
docker build -t registry.yourcompany.com/monitoring-agent:latest .
docker push registry.yourcompany.com/monitoring-agent:latest
```

---

## 2️⃣ Host the Install Script

### Option A: Static File Hosting (Vercel/Netlify)
```bash
# Create a public directory
mkdir public
cp install.sh public/install.sh

# Deploy to Vercel
vercel --prod

# Your URL will be: https://yourapp.vercel.app/install.sh
```

### Option B: GitHub Pages
```bash
# Create gh-pages branch
git checkout -b gh-pages
cp install.sh index.html
git add install.sh
git commit -m "Add install script"
git push origin gh-pages

# Enable GitHub Pages in repo settings
# URL: https://yourusername.github.io/monitoring-agent/install.sh
```

### Option C: Your Own Domain
```bash
# Upload to your server
scp install.sh user@yourserver.com:/var/www/html/install.sh

# Or use S3/CloudFront
aws s3 cp install.sh s3://yourbucket/install.sh --acl public-read
```

**Important:** Update the `AGENT_IMAGE` variable in `install.sh` to point to your registry.

---

## 3️⃣ Backend API Setup

Your backend needs to accept metrics from agents.

### Endpoint Requirements
```
POST /api/metrics
Authorization: Bearer <agent_api_key>
Content-Type: application/json

Body: {
  "server_id": "...",
  "timestamp": "...",
  "system": {...},
  "containers": [...]
}
```

### Example Express.js Handler
```javascript
app.post('/api/metrics', authenticateAgent, async (req, res) => {
  try {
    const { server_id, system, containers } = req.body;
    
    // Store in database
    await db.metrics.create({
      server_id,
      timestamp: new Date(),
      cpu: system.cpu_percent,
      memory: system.memory_percent,
      containers: containers.length
    });
    
    // Push to mobile clients via WebSocket/FCM
    await notifyMobileClients(server_id, req.body);
    
    res.json({ success: true });
  } catch (error) {
    console.error('Metrics error:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});
```

### Authentication Middleware
```javascript
function authenticateAgent(req, res, next) {
  const token = req.headers.authorization?.replace('Bearer ', '');
  
  if (!token || !isValidApiKey(token)) {
    return res.status(401).json({ error: 'Invalid API key' });
  }
  
  req.userId = getUserIdFromToken(token);
  next();
}
```

---

## 4️⃣ User Onboarding Flow

### Step 1: User signs up in mobile app
- Create account
- Generate unique API key
- Store in database

### Step 2: Show installation instructions
```
In your mobile app, show:

┌─────────────────────────────────┐
│  Add Your First Server          │
├─────────────────────────────────┤
│                                 │
│  Run this on your server:       │
│                                 │
│  curl -fsSL                     │
│  https://yourapp.com/install.sh │
│  | bash                         │
│                                 │
│  Your API Key:                  │
│  sk_live_abc123... [Copy]       │
│                                 │
└─────────────────────────────────┘
```

### Step 3: User runs command
```bash
# On their server
curl -fsSL https://yourapp.com/install.sh | bash
# Paste API key when prompted
```

### Step 4: Agent connects
- Agent starts sending metrics within 10 seconds
- Backend receives first payload
- Mobile app shows "Server Connected!" 🎉

---

## 5️⃣ Production Considerations

### Rate Limiting
Protect your backend from abuse:
```javascript
const rateLimit = require('express-rate-limit');

const metricsLimiter = rateLimit({
  windowMs: 1 * 60 * 1000, // 1 minute
  max: 10, // 10 requests per minute per IP
  message: 'Too many requests'
});

app.post('/api/metrics', metricsLimiter, ...);
```

### Monitoring the Agent
Monitor your monitoring agent:
- Track agent heartbeats
- Alert if agent stops reporting
- Monitor backend API errors
- Track API key usage

### Database Schema
```sql
CREATE TABLE metrics (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  server_id VARCHAR(255) NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  cpu_percent DECIMAL(5,2),
  memory_percent DECIMAL(5,2),
  disk_percent DECIMAL(5,2),
  container_count INTEGER,
  raw_payload JSONB,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_metrics_user_server ON metrics(user_id, server_id);
CREATE INDEX idx_metrics_timestamp ON metrics(timestamp DESC);
```

### Retention Policy
```sql
-- Delete metrics older than 30 days
DELETE FROM metrics 
WHERE timestamp < NOW() - INTERVAL '30 days';

-- Or archive to cold storage
INSERT INTO metrics_archive 
SELECT * FROM metrics 
WHERE timestamp < NOW() - INTERVAL '7 days';

DELETE FROM metrics 
WHERE timestamp < NOW() - INTERVAL '7 days';
```

---

## 6️⃣ Scaling Strategy

### Phase 1: MVP (1-100 servers)
- Single backend server
- PostgreSQL database
- Simple REST API
- **Cost: ~$20/month**

### Phase 2: Growth (100-1000 servers)
- Load balancer
- Multiple backend instances
- Redis for caching
- Message queue (RabbitMQ/Redis)
- **Cost: ~$100-200/month**

### Phase 3: Scale (1000+ servers)
- Time-series database (InfluxDB/TimescaleDB)
- Horizontal scaling
- CDN for install script
- Metrics aggregation
- **Cost: ~$500+/month**

---

## 7️⃣ Troubleshooting Production Issues

### Agent not connecting
```bash
# On the server, check logs
docker logs monitoring-agent

# Common issues:
# 1. Wrong API key → Regenerate in dashboard
# 2. Firewall blocking HTTPS → Check outbound 443
# 3. Backend down → Check backend status page
```

### High backend load
```bash
# Check metrics rate
SELECT COUNT(*) as requests_per_minute
FROM metrics 
WHERE timestamp > NOW() - INTERVAL '1 minute';

# Identify heavy senders
SELECT server_id, COUNT(*) as request_count
FROM metrics
WHERE timestamp > NOW() - INTERVAL '1 hour'
GROUP BY server_id
ORDER BY request_count DESC
LIMIT 10;
```

### Database growing too fast
```bash
# Check table size
SELECT pg_size_pretty(pg_total_relation_size('metrics'));

# Implement aggressive retention
DELETE FROM metrics WHERE timestamp < NOW() - INTERVAL '24 hours';
```

---

## 8️⃣ Security Hardening

### API Key Management
- Generate cryptographically secure keys
- Allow key rotation
- Implement key revocation
- Monitor for leaked keys

```javascript
// Generate secure API key
const crypto = require('crypto');
function generateApiKey() {
  return 'sk_live_' + crypto.randomBytes(32).toString('hex');
}
```

### Network Security
- Use HTTPS everywhere
- Validate SSL certificates
- Implement CORS properly
- Rate limit aggressively

### Agent Permissions
- Read-only Docker socket
- No root privileges
- Minimal host access
- Regular security updates

---

## 9️⃣ Monitoring & Alerting

### Key Metrics to Track
- Agent connection rate
- API error rate
- Database query performance
- Backend CPU/memory
- Failed authentication attempts

### Alert Conditions
```
Alert: High Error Rate
Condition: API error rate > 5% for 5 minutes
Action: Page on-call engineer

Alert: Agent Offline
Condition: No metrics from server for 2 minutes
Action: Notify user via mobile push

Alert: Database Slow
Condition: Query time > 1s for 5 minutes
Action: Notify DevOps team
```

---

## 🔟 Launch Checklist

- [ ] Agent builds and pushes successfully
- [ ] Install script is publicly accessible
- [ ] Backend API is running and tested
- [ ] Database is set up with proper indexes
- [ ] Mobile app can display incoming metrics
- [ ] API keys can be generated and revoked
- [ ] Error tracking is configured (Sentry/Rollbar)
- [ ] Monitoring is set up for backend
- [ ] Documentation is published
- [ ] Support email is configured
- [ ] Pricing page is ready (if applicable)

---

## 📈 Post-Launch

### Week 1
- Monitor error rates closely
- Collect user feedback
- Fix critical bugs immediately
- Update documentation based on questions

### Month 1
- Analyze usage patterns
- Identify performance bottlenecks
- Plan feature roadmap
- Start thinking about scaling

### Month 3
- Implement automatic scaling
- Add more metrics
- Improve mobile UX
- Consider paid tiers

---

## 🆘 Getting Help

If you encounter issues during deployment:

1. Check logs: `docker logs monitoring-agent`
2. Test manually: `curl -X POST ... ` to your backend
3. Verify network: `ping api.yourapp.com`
4. Check status page
5. Contact support

---

## ✅ Success Criteria

Your deployment is successful when:

- ✅ Users can install agent with one command
- ✅ Metrics appear in mobile app within 30 seconds
- ✅ API error rate < 0.1%
- ✅ Backend response time < 200ms p95
- ✅ 99.9% agent uptime
- ✅ Zero security incidents

---

**Ready to deploy? Let's go! 🚀**
