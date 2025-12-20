#!/bin/bash

# Or with API key: curl -fsSL https://yourapp.com/install.sh | API_KEY=your-key bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}  Monitoring Agent Installer${NC}"
echo -e "${GREEN}=====================================${NC}"
echo ""

# Configuration
AGENT_IMAGE="${AGENT_IMAGE:-yourregistry/monitoring-agent:latest}"
AGENT_NAME="${AGENT_NAME:-monitoring-agent}"
BACKEND_URL="${BACKEND_URL:-https://api.yourapp.com}"
API_KEY="${API_KEY:-}"
INTERVAL="${AGENT_INTERVAL:-10s}"
ENV="${AGENT_ENV:-production}"
LOG_LEVEL="${LOG_LEVEL:-info}"

# Check if running as root (warn but don't fail)
if [ "$EUID" -ne 0 ]; then 
    echo -e "${YELLOW}⚠ Warning: Not running as root. Some operations may require sudo.${NC}"
    SUDO="sudo"
else
    SUDO=""
fi

# Function to check command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if Docker is installed
echo "Checking Docker installation..."
if ! command_exists docker; then
    echo -e "${RED}✗ Error: Docker is not installed.${NC}"
    echo ""
    echo "Please install Docker first:"
    echo "  Ubuntu/Debian: curl -fsSL https://get.docker.com | sh"
    echo "  Or visit: https://docs.docker.com/engine/install/"
    exit 1
fi

echo -e "${GREEN}✓${NC} Docker is installed"

# Check if Docker daemon is running
echo "Checking Docker daemon..."
if ! $SUDO docker info >/dev/null 2>&1; then
    echo -e "${RED}✗ Error: Docker daemon is not running.${NC}"
    echo ""
    echo "Please start Docker:"
    echo "  sudo systemctl start docker"
    echo "  or: sudo service docker start"
    exit 1
fi

echo -e "${GREEN}✓${NC} Docker daemon is running"

# Check Docker permissions
if [ "$EUID" -ne 0 ] && ! docker ps >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠ Warning: Current user doesn't have Docker permissions${NC}"
    echo "Adding user to docker group (you may need to log out and back in)..."
    $SUDO usermod -aG docker "$USER" || true
    echo "Using sudo for Docker commands..."
    SUDO="sudo"
fi

# Get API key if not provided
if [ -z "$API_KEY" ]; then
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo "You need an API key to connect the agent to your account."
    echo "Get your API key from: ${BLUE}https://yourapp.com/settings/api-keys${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    read -p "Enter your API key: " API_KEY
    
    if [ -z "$API_KEY" ]; then
        echo -e "${RED}✗ Error: API key is required${NC}"
        exit 1
    fi
fi

# Validate API key format (basic check)
if [ ${#API_KEY} -lt 10 ]; then
    echo -e "${YELLOW}⚠ Warning: API key seems too short. Are you sure it's correct?${NC}"
    read -p "Continue anyway? (y/N): " confirm
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Stop and remove existing agent if present
if $SUDO docker ps -a --format '{{.Names}}' | grep -q "^${AGENT_NAME}$"; then
    echo -e "${YELLOW}Found existing agent, stopping and removing...${NC}"
    $SUDO docker stop $AGENT_NAME >/dev/null 2>&1 || true
    $SUDO docker rm $AGENT_NAME >/dev/null 2>&1 || true
    echo -e "${GREEN}✓${NC} Old agent removed"
fi

# Create persistent data directory
echo "Creating persistent storage..."
DATA_DIR="/var/lib/pulse-agent"
$SUDO mkdir -p "$DATA_DIR"
$SUDO chmod 755 "$DATA_DIR"
echo -e "${GREEN}✓${NC} Storage created at $DATA_DIR"

# Pull latest agent image
echo "Pulling latest agent image..."
if ! $SUDO docker pull $AGENT_IMAGE; then
    echo -e "${RED}✗ Error: Failed to pull agent image${NC}"
    echo "Image: $AGENT_IMAGE"
    echo ""
    echo "Possible solutions:"
    echo "  1. Check your internet connection"
    echo "  2. Verify the image name is correct"
    echo "  3. Check if you need to login: docker login"
    exit 1
fi

echo -e "${GREEN}✓${NC} Agent image pulled successfully"

# Get hostname
HOSTNAME=$(hostname)

# Start the agent
echo "Starting monitoring agent..."
if ! $SUDO docker run -d \
    --name $AGENT_NAME \
    --restart unless-stopped \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -v $DATA_DIR:/root/.pulse \
    -e AGENT_API_KEY="$API_KEY" \
    -e AGENT_BACKEND_URL="$BACKEND_URL" \
    -e AGENT_INTERVAL="$INTERVAL" \
    -e AGENT_ENV="$ENV" \
    -e LOG_LEVEL="$LOG_LEVEL" \
    $AGENT_IMAGE; then
    
    echo -e "${RED}✗ Error: Failed to start agent${NC}"
    echo "View logs with: $SUDO docker logs $AGENT_NAME"
    exit 1
fi

# Wait for agent to start
echo "Waiting for agent to initialize..."
sleep 5

# Check if agent is running
if ! $SUDO docker ps --format '{{.Names}}' | grep -q "^${AGENT_NAME}$"; then
    echo -e "${RED}✗ Error: Agent failed to start${NC}"
    echo ""
    echo "Recent logs:"
    $SUDO docker logs --tail 20 $AGENT_NAME
    echo ""
    echo "View full logs: $SUDO docker logs $AGENT_NAME"
    exit 1
fi

# Success message
echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}  ✓ Agent installed successfully!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}Server Details:${NC}"
echo "  • Hostname: $HOSTNAME"
echo "  • Backend: $BACKEND_URL"
echo "  • Environment: $ENV"
echo "  • Interval: $INTERVAL"
echo "  • Data Dir: $DATA_DIR"
echo ""
echo -e "${BLUE}Your server should appear in your dashboard within 30 seconds.${NC}"
echo ""
echo -e "${BLUE}Useful Commands:${NC}"
echo "  • View logs:      $SUDO docker logs -f $AGENT_NAME"
echo "  • Stop agent:     $SUDO docker stop $AGENT_NAME"
echo "  • Start agent:    $SUDO docker start $AGENT_NAME"
echo "  • Restart agent:  $SUDO docker restart $AGENT_NAME"
echo "  • Remove agent:   $SUDO docker stop $AGENT_NAME && $SUDO docker rm $AGENT_NAME"
echo "  • Update agent:   $SUDO docker pull $AGENT_IMAGE && $SUDO docker restart $AGENT_NAME"
echo ""
echo -e "${GREEN}Installation complete! 🎉${NC}"