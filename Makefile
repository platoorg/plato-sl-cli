.PHONY: build test test-integration install clean fmt lint

BINARY_NAME=platosl
MAIN_PATH=./cmd/platosl

build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

test:
	go test -v -cover ./...

test-integration:
	go test -v -tags=integration ./...

install:
	go install $(MAIN_PATH)

clean:
	rm -rf bin/
	go clean

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

.DEFAULT_GOAL := build
