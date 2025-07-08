# Simple Makefile for SwitchConfigSim

# Build both programs
build:
	go build -o switchctl cli/main.go
	go build -o api-server api/server.go

# Build just the CLI tool
build-cli:
	go build -o switchctl cli/main.go

# Build just the API server
build-api:
	go build -o api-server api/server.go

# Run the CLI tool (build it first)
run-cli: build-cli
	./switchctl

# Run the API server (build it first) 
run-api: build-api
	./api-server

# Clean up built files
clean:
	rm -f switchctl
	rm -f api-server

# Install Go dependencies
deps:
	go mod download


# Run tests
test:
	go test ./...

# Show help
help:
	@echo "Available commands:"
	@echo "  make build     - Build both programs"
	@echo "  make build-cli - Build CLI tool only"
	@echo "  make build-api - Build API server only"
	@echo "  make run-cli   - Build and run CLI tool"
	@echo "  make run-api   - Build and run API server"
	@echo "  make clean     - Remove built files"
	@echo "  make deps      - Download dependencies"
	@echo "  make test      - Run tests" 