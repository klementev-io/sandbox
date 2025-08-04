APP_NAME=sandbox

.PHONY: run lint docker-build stress-test install

run:
	@echo "Running $(APP_NAME)..."
	go run ./cmd/$(APP_NAME) -c ./configs/config.yaml

lint:
	golangci-lint run ./...

docker-build:
	docker build -t sandbox:latest .

docker-run:
	docker run --name sandbox -p 8080:80 sandbox:latest

STRESS_RATE=1000
STRESS_DURATION=60s

stress-test:
	echo "GET http://127.0.0.1:8080/api/v1/health" | vegeta attack -rate=$(STRESS_RATE) -duration=$(STRESS_DURATION) | vegeta report

GOLANGCI_VERSION=v2.3.0
VEGETA_VERSION=v12.12.0

install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go install github.com/tsenart/vegeta/v12@$(VEGETA_VERSION)