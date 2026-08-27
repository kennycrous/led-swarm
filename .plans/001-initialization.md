# Implementation Plan 001: Project Initialization & Infrastructure Setup

## 1. Overview & Objectives

This implementation plan establishes the foundational repository layout, Go backend engine structure, Svelte 5 frontend workspace, Wails desktop integration, and multi-stage Docker server setup for **LED Swarm Orchestrator**.

---

## 2. Decided Architecture & Repository Layout

```
led-swarm/
├── backend/                  # Go 1.22+ backend engine
│   ├── go.mod                # Go module definition
│   ├── main.go               # Entry point (CLI flags: --server, --port, Wails app)
│   ├── app.go                # Wails application bindings & API interface
│   ├── db.go                 # Pure-Go CGO-free SQLite database manager
│   ├── mdns.go               # mDNS zero-config WLED scanner (_wled._tcp)
│   ├── wled.go               # WLED device client (JSON REST, WebSockets, DDP UDP)
│   ├── hub.go                # Real-time WebSocket hub for UI streaming
│   └── server.go             # Embedded HTTP server for headless / Docker mode
├── frontend/                 # Svelte 5 + Vite + Tailwind CSS UI
│   ├── src/
│   │   ├── app.css           # Tailwind + Cyberpunk Glassmorphic CSS theme
│   │   ├── App.svelte        # Main shell layout with Cyberpunk glass navigation
│   │   └── lib/              # Svelte components & state stores
│   ├── package.json          # npm dependencies (svelte 5, vite, tailwind, lucide)
│   ├── vite.config.js        # Vite bundler config with dist build target
│   └── tailwind.config.js    # Cyberpunk color palette & glassmorphism utilities
├── Dockerfile                # Multi-stage Docker build (Node 20 -> Go 1.22 -> Alpine)
├── wails.json                # Wails v2 desktop configuration
├── .docs/
│   └── architecture.md       # Full architecture specification
└── AGENTS.md                 # Agent context & guidelines
```

---

## 3. Step-by-Step Implementation Tasks

### Task 1: Go Backend Module Setup (`backend/`)
- [ ] Initialize Go module `led-swarm/backend` inside `backend/`.
- [ ] Add pure Go dependencies to `go.mod`:
  - `github.com/glebarez/go-sqlite` (CGO-free SQLite)
  - `github.com/grandcat/zeroconf` (mDNS scanner)
  - `github.com/gorilla/websocket` (WebSocket hub)
  - `github.com/wailsapp/wails/v2` (Desktop framework)
- [ ] Implement `db.go`: SQLite initialization with pure Go driver creating tables for `devices`, `groups`, `canvas_placements`, and `scenes`.
- [ ] Implement `mdns.go`: Background worker skeleton scanning for `_wled._tcp`.
- [ ] Implement `hub.go`: Thread-safe WebSocket hub broadcasting live state updates.
- [ ] Implement `server.go`: HTTP server serving embedded frontend static files on port 8080 when running in `--server` mode.
- [ ] Implement `main.go`: Parse command-line flags (`--server`, `--port`, `--db`) and initialize either Wails Desktop GUI or Headless Web Server.

### Task 2: Svelte 5 Frontend Setup (`frontend/`)
- [ ] Initialize Svelte 5 + Vite application inside `frontend/`.
- [ ] Install dependencies in `frontend/package.json`:
  - `svelte@next` (Svelte 5)
  - `vite`
  - `tailwindcss`, `postcss`, `autoprefixer`, `@tailwindcss/vite`
  - `lucide-svelte`
  - `clsx`, `tailwind-merge`
- [ ] Configure `tailwind.config.js`:
  - Void background `#06090e`, cyber slate `#0c1017`, cyan neon `#06b6d4`, magenta `#a855f7`, gold `#f59e0b`.
  - Custom glassmorphism classes (`backdrop-blur-xl`, border glow effects).
- [ ] Create `src/app.css` with global dark background, custom scrollbars, and neon glow utility classes.
- [ ] Build `src/App.svelte` Cyberpunk shell (Header, Sidebar Rail, Glassmorphic Dashboard placeholder).

### Task 3: Wails Desktop Integration (`wails.json`)
- [ ] Create `wails.json` linking `backend/` entry point and `frontend/` build target.
- [ ] Configure Wails window parameters:
  - Width: 1280, Height: 800, Minimum Size: 900x600.
  - Title: "LED Swarm Orchestrator".
  - Background Color: `#06090e` (Dark Void).
- [ ] Add Go struct methods in `backend/app.go` for Wails TypeScript auto-generation.

### Task 4: Docker & Server Mode Packaging (`Dockerfile`)
- [ ] Create multi-stage `Dockerfile`:
  - **Stage 1 (Node 20 Alpine)**: Install frontend dependencies & run `npm run build` to generate `frontend/dist`.
  - **Stage 2 (Golang 1.22 Alpine)**: Compile static Go binary embedding `frontend/dist`.
  - **Stage 3 (Alpine Latest)**: Minimal runner container with `led-swarm-server` binary exposing port `8080`.

---

## 4. Verification & Testing Criteria

1. **Go Backend Compilation**:
   - `cd backend && go build -o led-swarm-server .` compiles cleanly with zero CGO dependencies (`CGO_ENABLED=0`).
2. **Frontend Vite Build**:
   - `cd frontend && npm run build` produces optimized static assets in `frontend/dist/`.
3. **Headless Server Mode**:
   - Running `./backend/led-swarm-server --server` starts the server on port `8080` and serves the Svelte UI at `http://localhost:8080`.
4. **Docker Container Build**:
   - `docker build -t led-swarm:latest .` produces a working container under 25MB.
5. **Wails Desktop Packaging**:
   - `wails build` compiles the native OS desktop executable.

---

*Plan stored in `.plans/001-initialization.md`.*
