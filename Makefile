APP_NAME=sandbox

.PHONY: run lint setup install

run:
	@echo "Running $(APP_NAME)..."
	go run ./cmd/$(APP_NAME) -c ./configs/config.yaml

setup:
	go mod tidy

GOLANGCI_VERSION=v2.3.0

install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)

lint:
	golangci-lint run ./...