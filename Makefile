# ===== IP Accounting Bot — Makefile =====
# Usage:
#   make help
#   make env            # create .env from .env.example (if missing)
#   make deps fmt vet
#   make test | test-race | cover
#   make build          # build both binaries
#   make run-bot        # start bot (loads .env if present)
#   make migrate        # run migrations (loads .env if present)
#   make clean

# --- Helper to load .env like a shell (handles quotes correctly) ---
# Used as:  @$(envsh); command ...
envsh = set -a; [ -f .env ] && . ./.env; set +a

# --- Paths / binaries ---
BIN_DIR     := bin
BOT_BIN     := $(BIN_DIR)/ip_bot
MIG_BIN     := $(BIN_DIR)/migrate
PKG_ALL     := ./...

# --- Go build flags (customize if needed) ---
GOFLAGS     ?=
LDFLAGS     ?=
GCFLAGS     ?=

# --- Default target ---
.DEFAULT_GOAL := help

# --- Help ---
.PHONY: help
help:
	@echo "Targets:"
	@echo "  env            - copy .env.example to .env (if missing)"
	@echo "  deps           - go mod tidy"
	@echo "  fmt            - gofmt -s -w"
	@echo "  vet            - go vet"
	@echo "  test           - go test ./..."
	@echo "  test-race      - go test -race ./..."
	@echo "  cover          - tests with coverage report"
	@echo "  build          - build bot and migrate binaries"
	@echo "  build-bot      - build only bot binary"
	@echo "  build-migrate  - build only migrate binary"
	@echo "  run-bot        - run bot (loads .env)"
	@echo "  migrate        - run migrations (loads .env)"
	@echo "  clean          - remove build artifacts"

# --- Env bootstrap ---
.PHONY: env
env:
	@if [ -f .env ]; then \
		echo ".env already exists (skipping)"; \
	elif [ -f .env.example ]; then \
		cp .env.example .env && echo "Created .env from .env.example"; \
	else \
		echo "No .env.example found"; exit 1; \
	fi

# --- Hygiene ---
.PHONY: deps fmt vet
deps:
	go mod tidy

fmt:
	gofmt -s -w .

vet:
	go vet $(PKG_ALL)

# --- Tests ---
.PHONY: test test-race cover
test:
	go test $(PKG_ALL)

test-race:
	go test -race $(PKG_ALL)

cover:
	go test -coverprofile=coverage.out $(PKG_ALL)
	@echo ""
	@echo "Coverage summary:"
	go tool cover -func=coverage.out | tail -n 1
	@echo "Open HTML report: go tool cover -html=coverage.out"

# --- Build ---
.PHONY: build build-bot build-migrate
build: build-bot build-migrate

build-bot:
	@mkdir -p $(BIN_DIR)
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -gcflags "$(GCFLAGS)" -o $(BOT_BIN) ./cmd/bot

build-migrate:
	@mkdir -p $(BIN_DIR)
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -gcflags "$(GCFLAGS)" -o $(MIG_BIN) ./cmd/migrate

# --- Run (loads .env if present) ---
.PHONY: run-bot migrate
run-bot: build-bot
	@$(envsh); $(BOT_BIN)

migrate: build-migrate
	@$(envsh); $(MIG_BIN)

# --- Clean ---
.PHONY: clean
clean:
	rm -rf $(BIN_DIR) coverage.out
