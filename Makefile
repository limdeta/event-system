APP_NAME=event-system
CMD_DIR=./cmd/$(APP_NAME)
INTERNAL_DIR=./internal
TESTS_DIR=./tests

# Main commands

.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

.PHONY: run
run:
	go run $(CMD_DIR)

.PHONY: install
install:
	go mod tidy
	go mod download

# Тесты

.PHONY: test
test: ## Unit + integration tests
	go test -count=1 -v $(INTERNAL_DIR)/...

.PHONY: test-e2e
test-e2e: ## End-to-end tests
	go test -count=1 -v $(TESTS_DIR)/e2e/...

.PHONY: test-bdd
test-bdd: ## BDD tests (requires godog)
	godog $(TESTS_DIR)/bdd/

.PHONY: test-all
test-all: test test-e2e ## All tests (except BDD)

# Docker

.PHONY: docker-build
docker-build:
	docker build -t $(APP_NAME):latest .

.PHONY: docker-run
docker-run:
	docker run --rm -it -p 8080:8080 $(APP_NAME):latest

.PHONY: compose-up
compose-up:
	docker-compose up -d

.PHONY: compose-down
compose-down:
	docker-compose down

# Code formatting and static analysis

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...

# Очистка

.PHONY: clean
clean:
	rm -rf bin

# help
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "  build         Build the project"
	@echo "  run           Run the application"
	@echo "  install       Install Go dependencies"
	@echo "  test          Run unit and integration tests"
	@echo "  test-e2e      Run end-to-end tests in ./tests/e2e/"
	@echo "  test-bdd      Run BDD tests in ./tests/bdd/ (needs godog installed)"
	@echo "  test-all      Run all tests except BDD"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-run    Run Docker container"
	@echo "  compose-up    Start docker-compose services"
	@echo "  compose-down  Stop docker-compose services"
	@echo "  fmt           Format Go code"
	@echo "  lint          Run linter (needs golangci-lint)"
	@echo "  clean         Remove build artifacts"