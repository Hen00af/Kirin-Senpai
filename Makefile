.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o discord-bot main.go

# Run the application (requires DISCORD_TOKEN environment variable)
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f discord-bot

# Build Docker image
docker-build:
	docker build -t discord-atcoder-bot .

# Run Docker container (requires DISCORD_TOKEN environment variable)
docker-run:
	docker run -e DISCORD_TOKEN=${DISCORD_TOKEN} discord-atcoder-bot

# Run with docker-compose
docker-compose-up:
	docker-compose up --build

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  docker-compose-up - Run with docker-compose"
	@echo "  deps            - Install dependencies"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code"
	@echo "  help            - Show this help"