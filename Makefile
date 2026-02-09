.PHONY: build test test-integration install clean fmt lint version

BINARY_NAME=platosl
MAIN_PATH=./cmd/platosl

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags with version information
LDFLAGS=-ldflags "-X github.com/platoorg/platosl-cli/internal/cli.Version=$(VERSION) \
                  -X github.com/platoorg/platosl-cli/internal/cli.Commit=$(COMMIT) \
                  -X github.com/platoorg/platosl-cli/internal/cli.BuildDate=$(BUILD_DATE)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PATH)

test:
	go test -v -cover ./...

test-integration:
	go test -v -tags=integration ./...

install:
	go install $(LDFLAGS) $(MAIN_PATH)

version:
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

clean:
	rm -rf bin/
	go clean

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

.DEFAULT_GOAL := build
