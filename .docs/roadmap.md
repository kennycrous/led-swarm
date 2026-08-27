# LED Swarm Orchestrator - Incremental Horizontal Feature Roadmap

## 1. Executive Summary & Horizontal Slicing Philosophy

This document defines the incremental, horizontally sliced feature roadmap for **LED Swarm Orchestrator**. 

### What is Horizontal Slicing?
Each phase in this roadmap delivers an **end-to-end working slice of functionality** — spanning backend Go logic, SQLite persistence, WLED network communication, and Svelte 5 Cyberpunk UI controls. The application remains **100% buildable, runnable, and testable at the end of every slice**.

---

## 2. Horizontal Feature Slices

```
+-----------------------------------------------------------------------------------+
|  Slice 001 [COMPLETED]: Infrastructure, Svelte UI Shell, Go Engine, Docker & Make |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 002 [COMPLETED]: Live WLED Device Discovery & Real-Time Controls           |
|  - Real mDNS & WLED WebSocket proxy live status updates                           |
|  - Power, brightness, color picker, effects, SQLite device storage               |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 003 [COMPLETED]: Virtual Strip Grouping & Multi-Zone Scenes                |
|  - Multi-strip groups ("Desk", "Ceiling", "TV Backlight")                         |
|  - Multi-device Scene JSON snapshots & 1-click restore                            |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 004 [COMPLETED]: Customizable Dashboard & Device Visibility Management     |
|  - Pin / unpin individual strips, groups, or scene presets to dashboard           |
|  - Settings -> Strips & Devices Management view showing all network strips        |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 005 [COMPLETED]: Automated Testing, Linting & Code Quality Suite           |
|  - Go & Vitest test suites, mock WLED server harness, GitHub Actions CI           |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 006 [COMPLETED]: Multi-Room 2D Layout Canvas Editor & Scoped Scenes        |
|  - Interactive 2D drag-and-drop physical strip room map                           |
|  - Multi-room canvas presets, dashboard room cards, scoped scene snapshots        |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 007 [COMPLETED]: Card-Integrated High-FPS DDP Pixel Streamer & 2D Spatial  |
|  - 60 FPS DDP (UDP 4048) target-keyed streamer (Device, Group, 2D Room Cards)     |
|  - 2D spatial coordinate generators (radial ripple, directional room sweep)       |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 008: Card-Integrated Sound Sync & Audio Reactivity Engine                  |
|  - Audio input FFT frequency binning (Bass, Mid, Treble)                          |
|  - Target-keyed sound reactivity toggles per Device, Group & 2D Room Canvas Card   |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 009: Desktop Integration, Settings & Production Release                    |
|  - System tray icon & minimize-to-tray in Wails                                   |
|  - Manual IP subnet scanner & GitHub Actions ghcr.io Docker workflow              |
+-----------------------------------------------------------------------------------+
```

---

## 3. Slice Breakdown & Acceptance Criteria

### 🟢 Slice 001 [COMPLETED]: Infrastructure & Baseline Foundation
- [x] Go 1.22+ backend engine structure (`backend/src`).
- [x] Pure-Go CGO-free SQLite driver integration (`glebarez/go-sqlite`).
- [x] Svelte 5 + Vite + Tailwind CSS Cyberpunk Glassmorphic UI shell (`frontend/`).
- [x] Wails v2 desktop configuration (`wails.json`).
- [x] Embedded HTTP server serving frontend static assets (`go:embed all:dist`).
- [x] Multi-stage `Dockerfile` (<25MB image) and development `Makefile`.
- [x] Air hot-reload integration (`backend/.air.toml`).

---

### 🟢 Slice 002 [COMPLETED]: Live WLED Device Discovery & Real-Time Controls
- **Objective**: Replace frontend sample placeholder data with live, auto-discovered WLED network devices and real-time controls.
- **Backend Deliverables**:
  - [x] WLED WebSocket proxy maintaining persistent connections to discovered WLED IPs.
  - [x] Broadcast live device state deltas to Svelte clients over `/api/v1/ws`.
  - [x] Save newly discovered devices and online/offline health updates into SQLite `devices` table.
  - [x] HTTP REST endpoints for custom device renaming, state updates, and manual IP addition.
  - [x] Cache built-in WLED effects (`/api/v1/effects`) and palettes (`/api/v1/palettes`).
- **Frontend Deliverables**:
  - [x] Connect Svelte UI to `/api/v1/ws` WebSocket channel.
  - [x] Real-time device cards displaying live WLED power status, MAC, IP, and LED count.
  - [x] Interactive controls: Power toggle, master brightness slider, color spectrum picker, effect & palette dropdowns.
  - [x] Device inline nickname editor for setting custom names in SQLite.
  - [x] Manual IP address addition modal for networks without mDNS.
- **Acceptance Criteria**:
  - [x] Booting backend on a network with WLED devices auto-populates the Svelte UI in real time.
  - [x] Toggling power, changing brightness, setting colors, or picking effects in Svelte UI instantly updates the physical WLED strip.

---

### 🟢 Slice 003 [COMPLETED]: Virtual Strip Grouping & Multi-Zone Scenes
- **Objective**: Allow users to combine physical WLED strips into logical groups and save multi-device snapshot scenes.
- **Backend Deliverables**:
  - [x] Database handlers for `groups`, `group_devices`, and `scenes` tables.
  - [x] Multi-device batch dispatcher applying state changes concurrently across grouped strips using `sync.WaitGroup`.
  - [x] Scene engine capturing multi-strip JSON state snapshots and restoring them in <50ms.
- **Frontend Deliverables**:
  - [x] Groups & Scenes tab in Svelte 5 UI (`GroupCard.svelte`, `SceneCard.svelte`).
  - [x] Group control cards for batch power toggle, group brightness scaling, color selection, and effect dropdowns.
  - [x] Scene Preset bar in top cyber header for 1-click scene activation anywhere in the application.
  - [x] Dashboard group filter pills for filtering device grid by active group.
- **Acceptance Criteria**:
  - [x] Creating a group "Desk Setup" containing multiple strips allows controlling all assigned strips simultaneously.
  - [x] Clicking a saved Scene button updates all assigned WLED strips concurrently in <50ms.

---

### 🟢 Slice 004 [COMPLETED]: Customizable Dashboard & Device Visibility Management
- **Objective**: Enable users to pin/unpin individual WLED light strips, virtual groups, or scene presets to/from the main Cyber Dashboard layout, and provide a dedicated Strips & Devices management view in Settings.
- **Backend Deliverables**:
  - [x] Database table `dashboard_items` (item_id, item_type: 'device'|'group'|'scene', position, size, is_pinned).
  - [x] Auto-pinning logic: Newly discovered WLED devices automatically pin to the dashboard layout.
  - [x] REST API & Wails desktop endpoints for pinning, unpinning, card sizing, and reordering dashboard cards (`/api/v1/dashboard/pin`, `/api/v1/dashboard/size`, `/api/v1/dashboard/reorder`).
- **Frontend Deliverables**:
  - **Customizable Main Dashboard**:
    - [x] Single unified canvas grid displaying pinned cards (mix of Individual Strip Cards, Group Cards, and 1-Click Scene Preset Cards).
    - [x] Card size toggle (Compact, Normal, Wide) and quick "Pin / Unpin" toggle action on cards.
  - **Settings -> Strips & Devices Management View**:
    - [x] Complete inventory view of ALL discovered and saved WLED light strips on the network regardless of dashboard visibility status.
    - [x] Per-strip management: Pin/unpin toggle, nickname editor, manual IP adder, discover streams trigger, and delete/forget device action.
- **Acceptance Criteria**:
  - [x] User can unpin individual strips from the dashboard while keeping their parent Group card visible on the dashboard.
  - [x] Settings -> "Strips & Devices" displays all WLED strips on the network regardless of dashboard visibility status.
  - [x] Newly discovered WLED devices auto-appear on the dashboard and in Settings.

---

### 🟢 Slice 005 [COMPLETED]: Automated Testing, Linting & Code Quality Suite
- **Objective**: Establish comprehensive automated unit testing, component testing, static code analysis, and CI linting for both the Go backend and Svelte 5 frontend.
- **Backend Deliverables**:
  - [x] Go unit and integration tests (`go test ./...`) covering SQLite database queries, device manager state cache, group/scene batch dispatching, and REST API handlers.
  - [x] Mock WLED WebSocket and mDNS server test harnesses for offline network testing.
  - [x] Static analysis and linting setup (`golangci-lint` or standard `go vet`).
- **Frontend Deliverables**:
  - [x] ESLint & `svelte-check` type checking pipeline for Svelte 5 components (`npm run lint`, `npm run check`).
  - [x] Vitest + Testing Library component test suite verifying store reactivity, card rendering, and user interactions.
- **CI / Build Infrastructure**:
  - [x] Makefile target additions (`make test`, `make lint`) for unified local verification.
  - [x] GitHub Actions CI workflow configuration executing tests and linters on pull requests.
- **Acceptance Criteria**:
  - [x] `make test` runs all backend and frontend unit tests with 100% pass rate.
  - [x] `make lint` verifies code style and type safety across Go backend and Svelte frontend with zero warnings.

---

### 🟢 Slice 006 [COMPLETED]: Multi-Room 2D Layout Canvas Editor & Scoped Scenes
- **Objective**: Create an interactive 2D visual workspace where users drag and position WLED light strips according to physical room layouts, pin 2D Room Canvases to the Dashboard, and capture room-scoped scenes that don't override untargeted rooms.
- **Backend Deliverables**:
  - [x] Database table `canvas_placements` (pos_x, pos_y, rotation, scale, geometry) & `CanvasManager` WebSocket engine.
  - [x] Database table `canvas_rooms` (`id`, `title`, `width`, `height`) for multi-room canvas preset management.
  - [x] Scoped Scenes schema (`scope_type`: 'global' | 'room' | 'group', `target_id`): Captures & restores state ONLY for assigned strips.
  - [x] Dashboard item support for `'room'` item type (`/api/v1/dashboard/pin`).
- **Frontend Deliverables**:
  - [x] Interactive 2D room map grid editor in Svelte 5 (`CanvasEditor.svelte`) with live pixel mirroring & spatial sweep animation.
  - [x] Multi-Room Canvas Selector (create, switch, rename, delete 2D Room Canvases).
  - [x] Pinnable 2D Room Canvas Cards on Dashboard grid displaying mini live 2D pixel preview & "Edit Layout" button.
  - [x] Group Creation Modal supporting Group Types (Standard Group, 2D Spatial Room, Scoped Scene).
  - [x] Scoped Scene capture modal allowing room-specific snapshot creation.
- **Acceptance Criteria**:
  - [x] User can create multiple 2D Room Canvases, pin a Room Canvas to the Dashboard, and open its 2D layout editor directly.
  - [x] Applying a Scoped Scene updates the assigned Room or Group without touching unassigned rooms.

---

### 🟢 Slice 007 [COMPLETED]: Card-Integrated High-FPS DDP / UDP Pixel Streamer
- **Objective**: Stream high-frequency (60 FPS) raw RGB pixel buffers over DDP UDP port 4048 for custom matrix animations and 2D room spatial visuals embedded directly into Device, Group, and 2D Room cards.
- **Backend Deliverables**:
  - [x] DDP UDP packet streamer engine (`backend/src/ddp_streamer.go`) emitting frame buffers on UDP port 4048.
  - [x] Multi-target streaming manager (`"device:<id>"`, `"group:<id>"`, `"room:<id>"`).
  - [x] Procedural generator engine (rainbow wave, digital rain, pulse beads, cyber fire, 2D radial ripple, 2D directional sweep).
- **Frontend Deliverables**:
  - [x] DDP 60 mode toggle button and live 60.0 FPS badge embedded in `DeviceCard.svelte`, `GroupCard.svelte`, and `RoomCard.svelte`.
  - [x] Dynamic CyberSelect dropdown swap to DDP procedural presets when DDP 60 mode is active.
  - [x] Reactive target-keyed `ddpStore.svelte.js` tracking active DDP streams over WebSockets.
- **Acceptance Criteria**:
  - [x] Activating DDP stream sends smooth 60 FPS pixel animations across single strips, virtual groups, and 2D room spatial maps without EEPROM wear.

---

### 🔴 Slice 008: Card-Integrated Sound Sync & Audio Reactivity Engine
- **Objective**: Capture real-time system audio input and map frequency spectrum energy (Bass, Mid, Treble) to per-card LED spatial animations directly on Device Cards, Group Cards, and 2D Room Canvas Cards (not application-wide global overrides).
- **Backend Deliverables**:
  - Native audio capture / WebAudio API bridge with FFT frequency analyzer binning audio into Bass, Mid, and Treble energy channels.
  - Target-keyed audio energy dispatcher mapping FFT reactivity to specific targets (`"device:<id>"`, `"group:<id>"`, `"room:<id>"`).
  - Procedural sound-reactive 1D & 2D spatial generators (*Bass Pulse*, *Spectrum Waterfall*, *Beat Ripple*, *VU Meter*, *Treble Sparkle*).
- **Frontend Deliverables**:
  - Embedded `Sound Sync` mode toggle button on `DeviceCard.svelte`, `GroupCard.svelte`, and `RoomCard.svelte`.
  - Per-card audio sensitivity gain slider, frequency band selector (Bass / Mid / Treble / Full Spectrum), and sound-reactive preset dropdown.
  - Mini FFT spectrum visualizer indicator bar directly on card headers when Sound Sync is active for that specific room or group.
- **Acceptance Criteria**:
  - Enabling Sound Sync on a specific Room or Group card drives real-time audio reactivity for that card's assigned strips while other unselected rooms/groups continue running their independent effects or DDP streams.

---

### ⚪ Slice 009: Desktop Integration, Settings & Production Release
- **Objective**: Implement native OS system desktop features, manual subnet network scanning, and automated production deployment workflows.
- **Backend / Frontend Deliverables**:
  - System tray icon & minimize-to-tray behavior in Wails desktop app.
  - Network Settings: Manual IP range subnet scanner for bulk discovering WLED strips on custom subnets.
  - GitHub Actions CI/CD release workflow compiling cross-platform standalone executables (Windows, macOS, Linux) and publishing Docker container images to `ghcr.io`.
- **Acceptance Criteria**:
  - Closing Wails desktop app minimizes to system tray; docker container builds automatically on GitHub release tags.

---

*Roadmap specification stored under `.docs/roadmap.md`.*
