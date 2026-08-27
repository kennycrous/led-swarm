# Implementation Plan - Slice 005: Automated Testing, Linting & Code Quality Suite

## Overview
Establish a comprehensive automated unit testing, component testing, static code analysis, and CI linting suite for both the Go backend engine and Svelte 5 frontend. This guarantees zero-regression code quality across WLED state orchestration, SQLite database persistence, cyber glassmorphic UI components, and pull request workflows.

---

## Technical Stack & Tooling Decisions

| Component | Tool / Library | Strategy / Rationale |
| :--- | :--- | :--- |
| **Backend Testing Framework** | Go standard `testing` + `net/http/httptest` | Zero external test dependencies, zero CGO, fast execution. |
| **Backend Database Isolation** | Pure-Go SQLite in-memory (`file::memory:?cache=shared`) | Isolated, sub-millisecond DB test fixtures without disk I/O. |
| **Backend Static Analysis** | `go vet` + `golangci-lint` | Go idiomatic code style, concurrency safety, and code quality enforcement. |
| **Frontend Test Runner** | Vitest + `@testing-library/svelte` + `jsdom` | Vite-native, high-performance Svelte 5 component & store testing. |
| **Frontend Code Quality** | ESLint 9 + `eslint-plugin-svelte` + `svelte-check` | Comprehensive Svelte 5 rune, a11y, and type checking pipeline. |
| **Automated CI/CD Workflows** | GitHub Actions (`.github/workflows/ci.yml`) | Enforces `make test` and `make lint` on all PRs to `develop` and `main`. |

---

## Step-by-Step Implementation Phases

### Phase 1: Backend Go Unit & Integration Test Suite (`backend/src/`)

- [ ] **`backend/src/db_test.go`**:
  - Test SQLite schema creation using in-memory database fixture.
  - Test `SaveDevice`, `GetDevices`, `UpdateDeviceName`, and online status toggles.
  - Test `SaveGroup`, `GetGroups`, `DeleteGroup`, `SaveScene`, `GetScenes`, and `DeleteScene`.
  - Test `dashboard_items` and `dashboard_panels` UPSERT behavior (`SaveDashboardPanel`, `GetDashboardPanels`, `UpdateDashboardItemPanel`, `UpdateDashboardItemSize`, `PinDashboardItem`).

- [ ] **`backend/src/server_test.go`**:
  - Use `httptest.NewServer` to test REST API HTTP endpoints:
    - `GET /api/v1/devices` & `POST /api/v1/devices/add`
    - `GET /api/v1/groups` & `POST /api/v1/groups`
    - `GET /api/v1/scenes` & `POST /api/v1/scenes/capture`
    - `GET /api/v1/dashboard/items` & `POST /api/v1/dashboard/panel` & `/api/v1/dashboard/panels`

- [ ] **`backend/src/dashboard_manager_test.go`**:
  - Test `PinItem` panel/size retention logic during simulated mDNS discovery.
  - Verify WebSocket broadcast payloads (`dashboard_pin_updated`, `dashboard_panel_updated`).

---

### Phase 2: Frontend Type Checking & ESLint Setup (`frontend/`)

- [ ] Install ESLint and Svelte plugin:
  ```bash
  cd frontend && npm install -D eslint eslint-plugin-svelte svelte-check typescript
  ```
- [ ] Create `frontend/eslint.config.js` with Svelte 5 rune rules and a11y checks.
- [ ] Add scripts to `frontend/package.json`:
  ```json
  "scripts": {
    "lint": "eslint .",
    "check": "svelte-check --tsconfig ./jsconfig.json"
  }
  ```
- [ ] Run `npm run lint` and `npm run check` to verify zero warnings across `App.svelte`, `CyberSelect.svelte`, `DeviceCard.svelte`, `GroupCard.svelte`, `SceneCard.svelte`, and `StripsManagement.svelte`.

---

### Phase 3: Frontend Vitest & Component Unit Testing (`frontend/`)

- [ ] Install Vitest and Svelte Testing Library:
  ```bash
  cd frontend && npm install -D vitest @testing-library/svelte jsdom @testing-library/jest-dom
  ```
- [ ] Create `frontend/vite.config.js` test configuration section.
- [ ] **`frontend/src/lib/stores/dashboardStore.test.js`**:
  - Test reactive state runes for adding panels, setting panel IDs, pinning items, and upserting WebSocket updates without wiping panel IDs.
- [ ] **`frontend/src/lib/stores/deviceStore.test.js`**:
  - Test device state caching and `onlineDevices` filter reactivity.
- [ ] **`frontend/src/lib/components/CyberSelect.test.js`**:
  - Test glassmorphic popover opening/closing, option selection, and keyboard navigation.

---

### Phase 4: Unified Makefile & GitHub Actions CI (`.github/workflows/`)

- [ ] Update `Makefile` targets:
  - `make test-backend`: `cd backend && go test -v ./src/...`
  - `make test-frontend`: `cd frontend && npm run test`
  - `make test`: Run `test-backend` and `test-frontend`.
  - `make lint-backend`: `cd backend && go vet ./src/...`
  - `make lint-frontend`: `cd frontend && npm run lint && npm run check`
  - `make lint`: Run `lint-backend` and `lint-frontend`.
- [ ] Create `.github/workflows/ci.yml`:
  - Run on `push` to `develop`, `main` and `pull_request`.
  - Steps:
    1. Set up Go 1.22+ and Node.js 20+.
    2. Install frontend dependencies (`cd frontend && npm ci`).
    3. Run `make lint`.
    4. Run `make test`.

---

## Acceptance Criteria
- `make test` executes all Go backend and Vitest frontend test suites with 100% pass rate.
- `make lint` verifies code style, type safety, and Svelte 5 a11y rules with zero errors.
- Feature branch `005-testing-and-linting` passes CI pipeline before merging into `develop`.
