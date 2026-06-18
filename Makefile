BINARY_NAME=recycleapp-ics

.PHONY: all build test lint fmt clean run help

all: lint test build

build:
	go build -v -o $(BINARY_NAME) ./cmd/recycleapp-ics

run: build
	./$(BINARY_NAME)

test:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist

help:
	@echo "Available targets:"
	@echo "  build   - Build the binary"
	@echo "  run     - Build and run the binary"
	@echo "  test    - Run tests"
	@echo "  lint    - Run golangci-lint"
	@echo "  fmt     - Format Go source code"
	@echo "  clean   - Remove binary and build artifacts"
	@echo "  help    - Display this help message"
