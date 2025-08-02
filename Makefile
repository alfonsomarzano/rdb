# RDB CLI Makefile

.PHONY: build test clean install help

# Build the RDB CLI tool
build:
	go build -o bin/rdb.exe main.go

# Build for multiple platforms
build-all: build-linux build-windows build-darwin

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/rdb-linux-amd64 main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/rdb-windows-amd64.exe main.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/rdb-darwin-amd64 main.go

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out

# Install the tool
install:
	go install

# Run the tool
run:
	go run main.go

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the RDB CLI tool"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install the tool"
	@echo "  run           - Run the tool"
	@echo "  help          - Show this help"

# Default target
all: deps build test 