# Slice 007 Plan: Card-Integrated High-FPS DDP 60 Streamer & 2D Spatial Engine

## 1. Executive Overview & Objectives

**Slice 007** embeds high-frequency **60 FPS DDP (Distributed Display Protocol)** streaming directly into **Device Cards**, **Group Cards**, and **2D Room Canvas Cards**, removing the need for a separate standalone tab.

### Key Architecture Decisions:
- **Card-Integrated DDP Toggle**: Every Device, Group, and Room Card includes a `DDP 60` mode toggle.
- **Dynamic Preset Switching**: Turning DDP 60 ON swaps standard WLED presets for DDP 60 FPS Procedural Generators (*Rainbow Wave*, *Digital Rain*, *Pulse Beads*, *Cyberpunk Flame*, *2D Spatial Ripple*, *2D Spatial Sweep*).
- **Target-Aware Streaming**:
  - **Single Device**: Streams to 1 strip IP.
  - **Virtual Group**: Streams across all group strips as a continuous 1D array buffer.
  - **2D Room Canvas**: Renders 2D spatial math using physical $(x, y)$ room placement coordinates from `CanvasManager`.

---

## 2. Technical Architecture & Task Breakdown

### 2.1 Backend Engine (`backend/src/ddp_streamer.go` & `server.go`)
- [ ] Extend `DDPStreamer` to manage concurrent active streams keyed by target key (`"device:<id>"`, `"group:<id>"`, `"room:<id>"`).
- [ ] Implement target resolution helpers:
  - Resolve Group strips via `GroupManager`.
  - Resolve Room Canvas placements via `CanvasManager` and compute 2D physical coordinates $(x_i, y_i)$ for each LED pixel.
- [ ] Implement 2D Spatial Procedural Generators:
  - `spatial_ripple`: Expanding 2D radial wave centered on room canvas.
  - `spatial_sweep`: 2D directional neon beam sweeping across physical room grid.
- [ ] Update REST API & Wails bindings:
  - `POST /api/v1/ddp/start` (`{ targetType, targetID, effect, speed, intensity }`)
  - `POST /api/v1/ddp/stop` (`{ targetType, targetID }`)
  - `GET /api/v1/ddp/status`
- [ ] Broadcast real-time `ddp_status` updates over WebSocket hub.

### 2.2 Frontend Svelte 5 UI (`ddpStore.js`, `DeviceCard.svelte`, `GroupCard.svelte`, `RoomCard.svelte`)
- [ ] Update `ddpStore.svelte.js` to track `activeStreams` map keyed by `${targetType}:${targetID}`.
- [ ] Embed `DDP 60` toggle button, DDP preset selector, and Speed/Intensity controls in:
  - `DeviceCard.svelte`
  - `GroupCard.svelte`
  - `RoomCard.svelte`
- [ ] Clean up `App.svelte`: remove standalone DDP Studio tab and unused sidebar navigation item.

### 2.3 Automated Testing (`ddp_streamer_test.go`, `ddpStore.test.js`, Component Tests)
- [ ] Backend tests for target resolution (device, group, room 2D spatial math).
- [ ] Frontend tests for card DDP mode toggle, preset selection, and store stream dispatching.
