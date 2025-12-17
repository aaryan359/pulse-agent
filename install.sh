#!/bin/bash

# Monitoring Agent Bootstrap Installer
# Usage: curl -fsSL https://yourapp.com/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}  Monitoring Agent Installer${NC}"
echo -e "${GREEN}=====================================${NC}"
echo ""

# Configuration
AGENT_IMAGE="yourregistry/monitoring-agent:latest"
AGENT_NAME="monitoring-agent"
BACKEND_URL="${BACKEND_URL:-https://api.yourapp.com}"
API_KEY="${API_KEY:-}"

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${YELLOW}Warning: Not running as root. Some operations may require sudo.${NC}"
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed.${NC}"
    echo "Please install Docker first: https://docs.docker.com/engine/install/"
    exit 1
fi

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    echo -e "${RED}Error: Docker daemon is not running.${NC}"
    echo "Please start Docker and try again."
    exit 1
fi

echo -e "${GREEN}✓${NC} Docker is installed and running"

# Get API key if not provided
if [ -z "$API_KEY" ]; then
    echo ""
    echo "You need an API key to connect the agent to your account."
    echo "Get your API key from: https://yourapp.com/settings/api-keys"
    echo ""
    read -p "Enter your API key: " API_KEY
    
    if [ -z "$API_KEY" ]; then
        echo -e "${RED}Error: API key is required${NC}"
        exit 1
    fi
fi

# Stop and remove existing agent if present
if docker ps -a --format '{{.Names}}' | grep -q "^${AGENT_NAME}$"; then
    echo -e "${YELLOW}Stopping existing agent...${NC}"
    docker stop $AGENT_NAME &> /dev/null || true
    docker rm $AGENT_NAME &> /dev/null || true
fi

# Pull latest agent image
echo "Pulling latest agent image..."
if ! docker pull $AGENT_IMAGE; then
    echo -e "${RED}Error: Failed to pull agent image${NC}"
    exit 1
fi

echo -e "${GREEN}✓${NC} Agent image pulled successfully"

# Get hostname
HOSTNAME=$(hostname)

# Start the agent
echo "Starting monitoring agent..."
docker run -d \
    --name $AGENT_NAME \
    --restart unless-stopped \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -e AGENT_API_KEY="$API_KEY" \
    -e AGENT_BACKEND_URL="$BACKEND_URL" \
    -e AGENT_SERVER_ID="$HOSTNAME" \
    -e AGENT_INTERVAL="10" \
    -e LOG_LEVEL="info" \
    $AGENT_IMAGE

# Wait for agent to start
sleep 3

# Check if agent is running
if docker ps --format '{{.Names}}' | grep -q "^${AGENT_NAME}$"; then
    echo ""
    echo -e "${GREEN}=====================================${NC}"
    echo -e "${GREEN}  ✓ Agent installed successfully!${NC}"
    echo -e "${GREEN}=====================================${NC}"
    echo ""
    echo "Server ID: $HOSTNAME"
    echo "Backend: $BACKEND_URL"
    echo ""
    echo "Your server should appear in your mobile app within 30 seconds."
    echo ""
    echo "Useful commands:"
    echo "  View logs:    docker logs -f $AGENT_NAME"
    echo "  Stop agent:   docker stop $AGENT_NAME"
    echo "  Start agent:  docker start $AGENT_NAME"
    echo "  Remove agent: docker stop $AGENT_NAME && docker rm $AGENT_NAME"
    echo ""
else
    echo -e "${RED}Error: Agent failed to start${NC}"
    echo "View logs with: docker logs $AGENT_NAME"
    exit 1
fi