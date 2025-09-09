
.PHONY: build test run lint docker-up docker-down help

# Go variables
BINARY_NAME=main

# Default command
help:
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@echo "  build       	Build the Go application"
	@echo "  test        	Run tests and show coverage"
	@echo "  run         	Run the application locally (requires docker-up)"
	@echo "  lint        	Run the golangci-lint linter"
	@echo "  docker-up   	Start all services with docker-compose"
	@echo "  docker-down 	Stop all services with docker-compose"

# Build the Go application
build:
	@echo "Building the application..."
	go build -o $(BINARY_NAME) ./main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v -cover ./...

# Run the application using docker-compose
run:
	@echo "Running the application with docker-compose..."
	docker-compose up app

# Run the linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Start all docker services
docker-up:
	@echo "Starting docker services..."
	docker-compose up -d

# Stop all docker services
docker-down:
	@echo "Stopping docker services..."
	docker-compose down
