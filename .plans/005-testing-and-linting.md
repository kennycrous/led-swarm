# Implementation Plan - Slice 005: Automated Testing, Linting & Code Quality Suite

## Overview
Establish a robust automated testing and static analysis suite for both the Go backend and Svelte 5 frontend. This ensures high code quality, prevents regressions as complex features are added, and enforces unified linting across the repository.

---

## Technical Stack & Tools

| Component | Tool / Library | Purpose |
| :--- | :--- | :--- |
| **Backend Testing** | Go standard `testing` package + SQLite test DB | Unit and integration tests for DB, HTTP API, and WebSocket proxy. |
| **Backend Mocking** | `net/http/httptest` + Mock WLED UDP/WS server | Offline network testing without physical hardware. |
| **Backend Linting** | `go vet` + `golangci-lint` | Go static analysis and style checks. |
| **Frontend Type Checking** | `svelte-check` | Svelte 5 rune and prop type verification. |
| **Frontend Linting** | ESLint + `eslint-plugin-svelte` | Syntax, formatting, and a11y rule enforcement. |
| **Frontend Unit Testing** | Vitest + `@testing-library/svelte` | Store reactivity and Svelte component unit testing. |
| **CI/CD Automation** | GitHub Actions + Makefile | Local (`make test`, `make lint`) and PR workflow validation. |

---

## Key Phases

### Phase 1: Backend Testing Infrastructure (`backend/`)
1. Create `backend/src/db_test.go` verifying SQLite schema initialization, device CRUD, group/scene persistence, and `dashboard_items` / `dashboard_panels` UPSERT operations.
2. Create `backend/src/server_test.go` using `httptest.NewServer` to test REST API endpoints (`/api/v1/devices`, `/api/v1/groups`, `/api/v1/scenes`, `/api/v1/dashboard/*`).
3. Create `backend/src/device_manager_test.go` with mock WLED HTTP/WebSocket server handlers verifying mDNS auto-pinning and state caching.

### Phase 2: Frontend Linting & Type Checking (`frontend/`)
1. Configure ESLint and `svelte-check` in `frontend/package.json`.
2. Add scripts: `"lint": "eslint ."` and `"check": "svelte-check --tsconfig ./tsconfig.json"`.
3. Verify zero lint errors and zero type warnings across all Svelte 5 components.

### Phase 3: Frontend Unit Testing Suite (`frontend/`)
1. Install and configure Vitest and Svelte Testing Library.
2. Create unit tests for reactive stores:
   - `dashboardStore.test.js`: Panel addition, pinning, size toggling, and UPSERT retention.
   - `deviceStore.test.js`: State caching and `onlineDevices` filter.
   - `groupStore.test.js`: Group state dispatching and scene snapshot capture/restore.
3. Create component render tests for `CyberSelect.svelte`, `DeviceCard.svelte`, `GroupCard.svelte`, and `SceneCard.svelte`.

### Phase 4: Makefile Targets & GitHub Actions CI (`.github/workflows/`)
1. Add `make test` target executing both `cd backend && go test ./...` and `cd frontend && npm run test`.
2. Add `make lint` target executing `cd backend && go vet ./...` and `cd frontend && npm run lint && npm run check`.
3. Create `.github/workflows/ci.yml` triggering `make lint` and `make test` on pull requests to `develop` and `main`.

---

## Verification Criteria
- `make test` runs all backend and frontend unit tests with 100% pass rate.
- `make lint` executes static analysis and type checks with zero warnings or errors.
