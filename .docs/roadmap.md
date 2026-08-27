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
|  Slice 002: Live WLED Device Discovery & Real-Time Controls                      |
|  - Real mDNS & WLED WebSocket proxy live status updates                           |
|  - Power, brightness, color picker, effects, SQLite device storage               |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 003: Virtual Strip Grouping & Multi-Zone Scenes                            |
|  - Multi-strip groups ("Desk", "Ceiling", "TV Backlight")                         |
|  - Multi-device Scene JSON snapshots & 1-click restore                            |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 004: 2D Visual Layout Canvas Editor                                        |
|  - Interactive 2D drag-and-drop physical strip room map                           |
|  - Live pixel mirroring preview & spatial sweep animations                        |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 005: High-FPS DDP / UDP Pixel Streamer                                     |
|  - 60 FPS DDP (UDP 4048) pixel buffer generator                                   |
|  - Go procedural matrix animations (rainbow wave, digital rain, pulse)            |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 006: Sound Sync & Audio Reactivity Engine                                  |
|  - Audio input FFT frequency binning (Bass, Mid, Treble)                          |
|  - Audio energy mapping to LED spatial pulses & live spectrum UI                  |
+-----------------------------------------------------------------------------------+
                                         │
                                         ▼
+-----------------------------------------------------------------------------------+
|  Slice 007: Desktop Integration, Settings & Production Release                    |
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

### 🔵 Slice 002: Live WLED Device Discovery & Real-Time Controls
- **Objective**: Replace frontend sample placeholder data with live, auto-discovered WLED network devices and real-time controls.
- **Backend Deliverables**:
  - WLED WebSocket proxy maintaining persistent connections to discovered WLED IPs.
  - Broadcast live device state deltas to Svelte clients over `/api/v1/ws`.
  - Save newly discovered devices and online/offline health updates into SQLite `devices` table.
  - HTTP REST endpoints for custom device renaming and manual IP addition.
- **Frontend Deliverables**:
  - Connect Svelte UI to `/api/v1/ws` WebSocket channel.
  - Real-time device cards displaying live WLED power status, MAC, IP, and LED count.
  - Interactive controls: Power toggle, master brightness slider, color spectrum picker, effect & palette dropdowns.
  - Device edit modal for setting custom nicknames.
- **Acceptance Criteria**:
  - Booting backend on a network with WLED devices auto-populates the Svelte UI in real time.
  - Toggling power or changing brightness in Svelte UI instantly updates the physical WLED strip.

---

### 🟣 Slice 003: Virtual Strip Grouping & Multi-Zone Scenes
- **Objective**: Allow users to combine physical WLED strips into logical groups and save multi-device snapshot scenes.
- **Backend Deliverables**:
  - Database handlers for `groups`, `group_devices`, and `scenes` tables.
  - Multi-device batch dispatcher applying state changes concurrently across grouped strips.
  - Scene engine capturing multi-strip JSON state snapshots and restoring them.
- **Frontend Deliverables**:
  - Group Management tab in Svelte UI (Create, Edit, Delete groups).
  - Group control cards for batch power toggle, group brightness scaling, and unified color selection.
  - Scene Preset bar in top header: 1-click scene activation buttons ("Cyberpunk Cyan", "Movie Night", "Warm Relax").
- **Acceptance Criteria**:
  - Creating a group "Desk Setup" containing 2 strips allows controlling both strips simultaneously.
  - Clicking a saved Scene button updates all assigned WLED strips in under 50ms.

---

### 🟡 Slice 004: 2D Visual Layout Canvas Editor
- **Objective**: Create an interactive 2D visual workspace where users drag and position WLED light strips according to their real-world physical room layout.
- **Backend Deliverables**:
  - API endpoints and SQLite persistence for `canvas_placements` (pos_x, pos_y, rotation, scale, length, geometry_type).
  - Spatial coordinate mapping engine for multi-strip 2D sweep effects.
- **Frontend Deliverables**:
  - Interactive HTML5 Canvas / WebGL 2D editor grid.
  - Drag-and-drop strip placement, rotation handles, length adjustment, and strip vs matrix geometry selection.
  - Live LED pixel mirroring preview (canvas pixels mirror physical WLED color output).
  - Spatial sweep animation trigger (rainbow wave sweeping across 2D canvas coordinates).
- **Acceptance Criteria**:
  - User can drag strips onto 2D canvas, save positions, and reload page with layout preserved.
  - Canvas pixels reflect real-time live LED colors of physical strips.

---

### 🟠 Slice 005: High-FPS DDP / UDP Pixel Streamer
- **Objective**: Stream high-frequency (60 FPS) raw RGB pixel buffers over DDP UDP port 4048 for custom matrix animations and swarm visuals.
- **Backend Deliverables**:
  - DDP UDP packet generator emitting frame buffers on UDP port 4048.
  - Procedural effect generator engine (matrix digital rain, rainbow wave, pulse beads, noise fire).
- **Frontend Deliverables**:
  - DDP Streaming control panel with effect dropdown, speed slider, and stream start/stop switch.
  - Status bar metrics: Real-time DDP frame rate (FPS) counter and UDP network latency graph.
- **Acceptance Criteria**:
  - Activating DDP stream sends smooth 60 FPS pixel animations across multiple WLED strips without EEPROM wear.

---

## 🔴 Slice 006: Sound Sync & Audio Reactivity Engine
- **Objective**: Capture live audio input and map frequency spectrum energy to WLED spatial animations.
- **Backend / Frontend Deliverables**:
  - Audio capture integration (native audio input or WebAudio FFT bridge).
  - FFT frequency analyzer binning audio into Bass, Mid, and Treble energy channels.
  - Audio Reactivity Studio tab in UI: Live FFT visualizer graph, gain sensitivity slider, and beat pulse mapping controls.
  - Audio-driven spatial pulse animations across the 2D visual layout canvas.
- **Acceptance Criteria**:
  - Playing music triggers real-time LED pulses and color shifts matching beat energy.

---

## ⚪ Slice 007: Desktop Integration, Settings & Production Release
- **Objective**: Add desktop OS integration features and automated production deployment workflows.
- **Backend / Frontend Deliverables**:
  - System tray icon & minimize-to-tray behavior in Wails desktop app.
  - Network subnet scanner for manual IP range searching.
  - GitHub Actions CI/CD release workflow compiling cross-platform binaries and pushing Docker images to `ghcr.io`.
- **Acceptance Criteria**:
  - Closing Wails desktop app minimizes to system tray; docker container builds automatically on GitHub tags.

---

*Roadmap specification stored under `.docs/roadmap.md`.*
