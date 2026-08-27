.PHONY: all dev-backend dev-frontend dev-desktop build build-frontend build-backend test test-backend test-frontend lint lint-backend lint-frontend fmt clean help

# Default target
all: build

## dev-backend: Run Go backend server locally (with Air live reload if installed)
dev-backend:
	@echo "==> Starting Go backend server in dev mode..."
	@if command -v air >/dev/null 2>&1; then \
		echo "==> Air live reload detected!"; \
		cd backend && air; \
	elif [ -f "$(HOME)/go/bin/air" ]; then \
		echo "==> Air live reload detected in ~/go/bin!"; \
		cd backend && $(HOME)/go/bin/air; \
	else \
		echo "==> [Notice] Air not found in PATH. Install with: go install github.com/air-verse/air@latest"; \
		echo "==> Running with plain go run..."; \
		cd backend && go run ./src --server --port=8080; \
	fi

## dev-frontend: Run Svelte 5 + Vite UI dev server with hot reload (http://localhost:5173)
dev-frontend:
	@echo "==> Starting Svelte 5 Vite dev server with Hot Module Replacement..."
	cd frontend && npm run dev

## dev-desktop: Run Wails desktop live reload application
dev-desktop:
	@echo "==> Starting Wails desktop dev application..."
	wails dev

## build-frontend: Build frontend static assets into backend/src/dist
build-frontend:
	@echo "==> Building frontend Svelte web assets..."
	cd frontend && npm run build

## build-backend: Build standalone Go server binary
build-backend: build-frontend
	@echo "==> Building static Go server binary..."
	cd backend && CGO_ENABLED=0 go build -ldflags="-s -w" -o led-swarm-server ./src

## build: Full build (frontend assets + backend static server binary)
build: build-backend

## fmt: Automatically format Go backend and Svelte/JS frontend files (removes trailing whitespace)
fmt:
	@echo "==> Formatting Go backend code (go fmt)..."
	cd backend && go fmt ./src/...
	@echo "==> Formatting Svelte and JS frontend code (Prettier)..."
	cd frontend && npm run format

## test-backend: Run backend Go unit tests (ensures frontend assets built first)
test-backend: build-frontend
	@echo "==> Running Go backend unit tests..."
	cd backend && go test -v ./src/...

## test-frontend: Run frontend Vitest unit tests
test-frontend:
	@echo "==> Running frontend Vitest unit tests..."
	cd frontend && npm run test

## test: Run all backend and frontend unit tests
test: test-backend test-frontend

## lint-backend: Run Go static code analysis (go vet, ensures frontend assets built first)
lint-backend: build-frontend
	@echo "==> Running Go static analysis (go vet)..."
	cd backend && go vet ./src/...

## lint-frontend: Run frontend formatting check (Prettier), ESLint, and svelte-check
lint-frontend:
	@echo "==> Running frontend formatting check (Prettier), ESLint, and svelte-check..."
	cd frontend && npm run format:check && npm run lint && npm run check

## lint: Run all backend and frontend linters
lint: lint-backend lint-frontend

## clean: Clean generated binaries and build artifacts
clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf backend/led-swarm-server backend/led-swarm.exe backend/src/dist/assets backend/*.db *.db bin/ build/ frontend/dist

## help: Display available targets
help:
	@echo "LED Swarm Orchestrator - Development & Build Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make fmt            Auto-format all backend Go code and frontend Svelte/JS files"
	@echo "  make dev-backend    Start Go backend server locally (http://localhost:8080)"
	@echo "  make dev-frontend   Start Svelte 5 + Vite frontend dev server (http://localhost:5173)"
	@echo "  make dev-desktop    Start Wails desktop live reload application"
	@echo "  make build          Build production frontend assets and Go server binary"
	@echo "  make test           Run Go backend unit tests and Vitest frontend tests"
	@echo "  make lint           Run Go vet static analysis, Prettier check, ESLint, and svelte-check"
	@echo "  make clean          Remove build artifacts and test databases"
