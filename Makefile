.PHONY: run lint test gen-api stress-test docker-build docker-run install

APP_NAME=sandbox

run:
	@echo "Running $(APP_NAME)..."
	go run ./cmd/$(APP_NAME) --c ./configs/config.yaml

lint:
	@echo "Linting $(APP_NAME)..."
	go tool golangci-lint run ./...

test:
	@echo "Running tests $(APP_NAME)..."
	go test ./...

generate:
	go generate ./...

STRESS_RATE=1000
STRESS_DURATION=60s

stress-test:
	echo "GET http://127.0.0.1:8080/health" | go tool vegeta attack -rate=$(STRESS_RATE) -duration=$(STRESS_DURATION) | go tool vegeta report

docker-build:
	docker build -t sandbox:latest .

docker-run:
	docker run --rm \
		-p 8080:8080 \
		-v ./configs/config.yaml:/etc/sandbox/configs/config.yaml:ro \
		sandbox:latest --c /etc/sandbox/configs/config.yaml

GOLANGCI_VERSION=v2.3.1
VEGETA_VERSION=v12.12.0
OAPI_CODEGEN_VERSION=v2.5.0

install:
	go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go get -tool github.com/tsenart/vegeta/v12@$(VEGETA_VERSION)
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)