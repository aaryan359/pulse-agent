# Makefile for Monitoring Agent

.PHONY: help build run test clean docker-build docker-run install

# Variables
BINARY_NAME=agent
DOCKER_IMAGE=monitoring-agent
DOCKER_TAG=latest
REGISTRY=yourregistry

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the agent binary
	@echo "Building agent..."
	go build -o $(BINARY_NAME) cmd/agent/main.go
	@echo "✓ Build complete: ./$(BINARY_NAME)"

run: ## Run the agent locally
	@echo "Running agent..."
	@if [ -z "$$AGENT_API_KEY" ]; then \
		echo "Error: AGENT_API_KEY not set"; \
		echo "Usage: AGENT_API_KEY=your_key make run"; \
		exit 1; \
	fi
	go run cmd/agent/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies downloaded"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "✓ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run: ## Run agent in Docker
	@echo "Running agent in Docker..."
	@if [ -z "$$AGENT_API_KEY" ]; then \
		echo "Error: AGENT_API_KEY not set"; \
		echo "Usage: AGENT_API_KEY=your_key make docker-run"; \
		exit 1; \
	fi
	docker run --rm \
		--name $(DOCKER_IMAGE) \
		-v /var/run/docker.sock:/var/run/docker.sock:ro \
		-e AGENT_API_KEY=$$AGENT_API_KEY \
		-e AGENT_BACKEND_URL=$${AGENT_BACKEND_URL:-http://localhost:8000} \
		-e LOG_LEVEL=$${LOG_LEVEL:-debug} \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-push: docker-build ## Build and push Docker image to registry
	@echo "Pushing to registry..."
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "✓ Pushed to $(REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)"

install: build ## Install agent binary to /usr/local/bin
	@echo "Installing agent..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Installed to /usr/local/bin/$(BINARY_NAME)"

dev: ## Run in development mode with auto-reload (requires air)
	@if ! command -v air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@echo "Starting development server with hot reload..."
	air

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

lint: ## Run linter (requires golangci-lint)
	@if ! command -v golangci-lint > /dev/null; then \
		echo "Error: golangci-lint not installed"; \
		echo "Install: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi
	@echo "Running linter..."
	golangci-lint run
	@echo "✓ Lint complete"

# Example usage targets
example-local: ## Example: Run locally with test config
	@echo "Running example with test configuration..."
	AGENT_API_KEY=test-key-123 \
	AGENT_BACKEND_URL=http://localhost:8000 \
	AGENT_INTERVAL=5 \
	LOG_LEVEL=debug \
	go run cmd/agent/main.go

example-docker: ## Example: Run in Docker with test config
	@echo "Running Docker example..."
	docker run --rm \
		--name monitoring-agent-example \
		-v /var/run/docker.sock:/var/run/docker.sock:ro \
		-e AGENT_API_KEY=test-key-123 \
		-e AGENT_BACKEND_URL=http://host.docker.internal:8000 \
		-e AGENT_INTERVAL=5 \
		-e LOG_LEVEL=debug \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

.DEFAULT_GOAL := help