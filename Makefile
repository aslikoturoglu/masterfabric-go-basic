.PHONY: all build run generate generate-dart tidy lint test docker-up docker-down clean help

## ── Variables ─────────────────────────────────────────────────────────────────
BINARY      := bin/server
CLI_BINARY  := bin/masterfabric_go
CMD         := ./cmd/server
CLI_CMD     := ./cmd/masterfabric_go
COMPOSE     := docker compose -f deployments/docker-compose.yml

## ── Defaults ──────────────────────────────────────────────────────────────────
all: build

## build: compile the server binary
build:
	@echo "==> Building $(BINARY)"
	@CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BINARY) $(CMD)

## build-cli: compile the masterfabric_go CLI binary
build-cli:
	@echo "==> Building $(CLI_BINARY)"
	@CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(CLI_BINARY) $(CLI_CMD)

## run: build and run the server locally (requires infra running)
run: build
	@echo "==> Running $(BINARY)"
	@./$(BINARY)

## dev: run with air hot-reload (requires: go install github.com/air-verse/air@latest)
dev:
	@air

## generate: re-run gqlgen code generation
generate:
	@echo "==> Running gqlgen"
	@go run github.com/99designs/gqlgen generate --config gqlgen.yml

## generate-dart: generate the sdk/dart_go_api Dart package from GraphQL schema
generate-dart: build-cli
	@echo "==> Generating Dart SDK package"
	@$(CLI_BINARY) generate dart

## tidy: tidy go modules
tidy:
	@go mod tidy

## lint: run golangci-lint (requires: brew install golangci-lint)
lint:
	@golangci-lint run ./...

## test: run all tests
test:
	@go test -race -count=1 ./...

## docker-up: start all infrastructure + app containers
docker-up:
	@$(COMPOSE) up --build -d

## docker-infra: start only postgres, redis, rabbitmq (no app container)
docker-infra:
	@$(COMPOSE) up -d postgres redis rabbitmq

## docker-down: stop and remove containers
docker-down:
	@$(COMPOSE) down

## docker-logs: tail app container logs
docker-logs:
	@$(COMPOSE) logs -f app

## clean: remove built binaries
clean:
	@rm -rf bin/

## help: print this message
help:
	@grep -E '^## ' Makefile | sed 's/^## //'
