# Payara Go SDK — Makefile
# Use: make [target]

.PHONY: build test test-integration lint clean run-payment run-withdrawal help

# Go
GO := go
GOFLAGS := -v
MODULE := github.com/payara-id/go-sdk

# Binaries (for examples)
BIN_DIR := bin

help:
	@echo "Targets:"
	@echo "  build            Build all packages"
	@echo "  test             Run unit tests"
	@echo "  test-integration Run integration tests (needs PAYARA_APP_ID, PAYARA_APP_SECRET)"
	@echo "  lint             Run golangci-lint (if installed)"
	@echo "  clean            Remove build artifacts and binaries"
	@echo "  run-payment      Build and run example payment_service"
	@echo "  run-withdrawal   Build and run example withdrawal_service"

build:
	$(GO) build $(GOFLAGS) ./...

test:
	$(GO) test -count=1 ./payara/...

test-integration:
	$(GO) test -tags=integration -count=1 $(GOFLAGS) ./payara/...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || (echo "golangci-lint not installed; run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

clean:
	$(GO) clean -testcache -cache
	rm -rf $(BIN_DIR)
	rm -rf coverage.out coverage.html

run-payment: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/payment_service ./example/payment_service
	@echo "Run with: PAYARA_APP_ID=... PAYARA_APP_SECRET=... ./$(BIN_DIR)/payment_service"
	./$(BIN_DIR)/payment_service

run-withdrawal: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/withdrawal_service ./example/withdrawal_service
	@echo "Run with: PAYARA_APP_ID=... PAYARA_APP_SECRET=... ./$(BIN_DIR)/withdrawal_service"
	./$(BIN_DIR)/withdrawal_service

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Optional: coverage
coverage:
	$(GO) test -coverprofile=coverage.out ./payara/...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in a browser"
