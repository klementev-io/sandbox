include ./configs/.env
export

GOLANGCI_VERSION=v2.3.0

.PHONY: run lint setup install

run:
	@echo "Running sandbox..."
	go run ./cmd/sandbox -c ./configs/config.yaml

lint:
	golangci-lint run ./...

install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)