# LED Swarm Orchestrator - Architecture & System Design

## 1. Executive Summary & Vision

**LED Swarm Orchestrator** is a high-performance, lightweight, cross-platform application and server designed to manage, orchestrate, and synchronize multiple WLED light strips seamlessly. Standard solutions suffer from outdated dependencies, bloated desktop runtimes, or lack of multi-device synchronization capabilities. 

LED Swarm bridges this gap by combining:
- A **Go (Golang)** backend engine providing ultra-fast mDNS device discovery, realtime WebSocket status aggregation, low-latency DDP/UDP pixel streaming, sound reactivity, and persistent SQLite storage.
- A **Svelte 5 + Tailwind CSS** frontend delivering a responsive, dark glassmorphic, cyberpunk aesthetic with real-time canvas visualizations and fluid animations.
- Dual-mode deployment: **Native Desktop Application** (via Wails v2/v3 using native OS webview) and **Headless Server Mode** (single zero-dependency Go static binary with embedded web assets, ideal for Linux, Docker, or Raspberry Pi).

---

## 2. Architectural Overview

```
                          +-----------------------------------+
                          |     Native Desktop App (Wails)    |
                          |  OR Headless Web Client Browser   |
                          +-----------------+-----------------+
                                            |
                                  Cyberpunk Web UI
                             (Svelte 5 + Tailwind CSS)
                                            |
                         HTTP REST API & WebSocket Channel
                                            v
+-----------------------------------------------------------------------------------+
|                            LED Swarm Go Backend Engine                            |
|                                                                                   |
|  +-------------------+  +--------------------+  +------------------------------+  |
|  |  mDNS Discovery   |  |   WLED Sync Hub    |  |  2D Visual Canvas Engine     |  |
|  |  (_wled._tcp)     |  | (JSON API / WS)    |  | (Coordinates, Zone Mapping)  |  |
|  +-------------------+  +--------------------+  +------------------------------+  |
|                                                                                   |
|  +-------------------+  +--------------------+  +------------------------------+  |
|  | UDP / DDP Stream  |  |  Audio Reactivity  |  |      SQLite Persistent       |  |
|  | (Pixel Streaming) |  |   Sound Sync Engine|  |      State Storage           |  |
|  +-------------------+  +--------------------+  +------------------------------+  |
+-----------------------------------------------------------------------------------+
       |                        |                         |
       | DDP/UDP Realtime       | HTTP / WS JSON          | mDNS Discovery
       v                        v                         v
+-----------------------------------------------------------------------------------+
|                              Target WLED Light Strips                             |
|          [ Strip 1 ]               [ Strip 2 ]              [ Strip N... ]        |
+-----------------------------------------------------------------------------------+
```

---

## 3. Technology Stack

### 3.1 Backend Engine (Go)
- **Language**: Go 1.22+
- **Persistence**: SQLite via CGO-free pure Go driver (`modernc.org/sqlite` / `github.com/glebarez/go-sqlite`) for seamless cross-compilation without C toolchains.
- **Discovery**: Zero-configuration mDNS scanning via `github.com/grandcat/zeroconf` targeting service `_wled._tcp`.
- **WLED Protocols**:
  - **JSON REST API / WebSockets**: State polling, preset execution, brightness, palette switching, segment control.
  - **DDP (Distributed Display Protocol)**: High frame rate (60 FPS) raw RGB/RGBW packet streaming over UDP port `4048`.
  - **WLED Realtime UDP Protocol**: Direct sync packet broadcast on UDP port `21324`.
- **Audio Processing**: Native audio capture + FFT frequency binning using Go audio bindings (`github.com/gordonklaus/portaudio` or WebAudio browser fallback bridge).

### 3.2 Frontend (Svelte 5 + Tailwind CSS)
- **Framework**: Svelte 5 (Runes `$state`, `$derived`, `$effect` for minimal overhead and instant reactivity).
- **Build Tool**: Vite.
- **Styling**: Tailwind CSS with custom glassmorphism plugins (`backdrop-blur-xl`, border glow effects, dark void slate palette).
- **Icons**: `lucide-svelte`.
- **Canvas / Layout Engine**: HTML5 2D Canvas / WebGL for rendering interactive physical strip locations, live pixel colors, drag-and-drop zone boundaries, and real-time visualizer.

### 3.3 Packaging & Distribution
- **Desktop Packaging**: Wails v2/v3 (Go + OS Native WebView: WebKit on macOS, WebView2 on Windows, WebKitGTK on Linux). Zero Electron bloat; final binary size ~15MB.
- **Headless Server Packaging**: Single standalone binary compiled with Go `go:embed` embedding `dist/` web assets.
- **Containerization**: Multi-stage Dockerfile producing a minimal Scratch or Alpine-based container (<25MB image size).

---

## 4. Key Subsystems & Core Capabilities

### 4.1 Auto-Discovery & Device Lifecycle Manager
- Background worker continuously monitors local network for `_wled._tcp` mDNS service broadcasts.
- Automatic IP resolution, hostname mapping, MAC address fingerprinting, and hardware info query (`/json/info`).
- Fallback manual IP address scanning & range discovery.
- Offline/Online health check heartbeat with exponential backoff retry.

### 4.2 Real-Time State Orchestration Hub
- Maintains a bidirectional state cache of all connected WLED strips.
- Synchronizes changes instantly across single devices, virtual groups, or all devices in unison.
- WebSocket broadcast engine propagates status changes to all connected UI clients within <5ms.

### 4.3 Virtual Strip Grouping & Multi-Zone Scenes
- Allows arbitrary grouping of physical strips (e.g., "Desk Setup", "Living Room Ceiling", "TV Backlight").
- Synchronized preset triggers, master brightness scaling, color palette matching, and unified effect speeds.
- Multi-zone scene management: Save multi-strip snapshots as custom Scenes and restore them with one click.

### 4.4 2D Visual Canvas Engine
- Interactive 2D workspace where users position LED strips according to their real-world physical layout.
- Supports straight strips, matrices, curves, and custom segment paths.
- Real-time live mirroring: Canvas pixels mirror exact live LED output of connected WLED strips.
- Spatial effect mapping: Apply spatial sweep effects across multiple strips based on 2D coordinates.

### 4.5 Low-Latency DDP/UDP Pixel Streamer
- High-frequency (30–60 FPS) direct pixel buffer generator.
- Bypasses WLED EEPROM wear by streaming volatile frame buffers over DDP UDP packets.
- Enables custom Go-based animations, matrix video streaming, and reactive color sweeps across multi-strip swarms.

### 4.6 Audio Reactivity & Sound Sync Engine
- Processes incoming audio streams into real-time frequency bands (Bass, Mid, Treble).
- Maps audio energy levels dynamically to strip brightness, color spectrum shifts, or spatial pulse waves on the 2D layout canvas.

---

## 5. Cyberpunk Glassmorphic UI/UX Specification

### 5.1 Color Palette & Theme Tokens
- **Background Slate**: `#06090e` (Deep Void) / `#0c1017` (Dark Cyber Slate)
- **Glass Panel Surface**: `rgba(15, 23, 42, 0.65)` with `backdrop-filter: blur(16px)`
- **Border Highlights**: `rgba(56, 189, 248, 0.2)` with glowing hover states `rgba(168, 85, 247, 0.4)`
- **Accent Neon Palette**:
  - **Cyan Neon**: `#06b6d4` / `#22d3ee`
  - **Purple/Magenta**: `#a855f7` / `#ec4899`
  - **Cyber Amber/Yellow**: `#f59e0b` / `#eab308`
  - **Status Emerald**: `#10b981` (Online), `#ef4444` (Offline)

### 5.2 UI Layout Structure
1. **Top Cyber Bar**: Logo badge, quick swarm power switch, global brightness slider, active swarm scene selector, network status pill.
2. **Left Navigation Rail**: Collapsible icons (Dashboard, Devices, Canvas Layout, Groups & Scenes, Audio Reactive, Settings).
3. **Main Content View**:
   - **Dashboard**: Card grid of WLED strips with live preview canvas thumbnails, color pickers, effect dropdowns.
   - **Visual Canvas Editor**: Interactive 2D drag-and-drop spatial room layout view.
   - **Audio Reactivity Studio**: Live FFT visualizer and gain controls.
4. **Bottom Dock / Status Bar**: Real-time DDP streaming FPS indicator, latency graph, total active LEDs count.

---

## 6. Database Schema (SQLite)

```sql
-- Registered WLED Devices
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,               -- MAC Address or UUID
    name TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    mac_address TEXT UNIQUE,
    led_count INTEGER DEFAULT 0,
    is_online BOOLEAN DEFAULT FALSE,
    last_seen DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Virtual Strip Groups
CREATE TABLE IF NOT EXISTS groups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Group Memberships
CREATE TABLE IF NOT EXISTS group_devices (
    group_id TEXT NOT NULL,
    device_id TEXT NOT NULL,
    PRIMARY KEY (group_id, device_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

-- 2D Canvas Rooms / Presets
CREATE TABLE IF NOT EXISTS canvas_rooms (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    width INTEGER DEFAULT 2000,
    height INTEGER DEFAULT 1200,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2D Visual Layout Canvas Placements
CREATE TABLE IF NOT EXISTS canvas_placements (
    device_id TEXT NOT NULL,
    room_id TEXT NOT NULL DEFAULT 'default',
    pos_x REAL NOT NULL DEFAULT 100.0,
    pos_y REAL NOT NULL DEFAULT 100.0,
    rotation REAL NOT NULL DEFAULT 0.0,
    scale REAL NOT NULL DEFAULT 1.0,
    geometry TEXT DEFAULT 'strip',
    PRIMARY KEY (device_id, room_id),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

-- Saved Swarm Scenes (Global or Room/Group Scoped)
CREATE TABLE IF NOT EXISTS scenes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    icon TEXT,
    scope_type TEXT DEFAULT 'global',   -- 'global', 'room', 'group'
    target_id TEXT DEFAULT '',          -- Room ID or Group ID if scoped
    config_json TEXT NOT NULL,          -- Targeted multi-device JSON snapshot
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 7. API & WebSocket Specification

### REST API Endpoints
- `GET /api/v1/devices` - List all discovered/configured devices.
- `POST /api/v1/devices/scan` - Trigger mDNS network scan.
- `POST /api/v1/devices/:id/state` - Set state (power, brightness, color, effect, segment).
- `GET /api/v1/groups` - Get virtual strip groups.
- `POST /api/v1/groups` - Create virtual strip group.
- `POST /api/v1/groups/:id/state` - Send batch command to group.
- `GET /api/v1/canvas` - Get 2D layout coordinates.
- `POST /api/v1/canvas` - Save 2D layout placements.
- `GET /api/v1/scenes` - Get saved scenes.
- `POST /api/v1/scenes/apply/:id` - Trigger scene activation across swarm.

### WebSocket Endpoint
- `WS /api/v1/ws` - Unified realtime channel:
  - Inbound: State change commands, live canvas cursor updates.
  - Outbound: Real-time device state pushes, mDNS device join/leave events, audio energy spectrum levels, DDP frame rate metrics.

---

## 8. Build & Packaging Matrix

| Target Platform | Packaging Tool | Output |
| :--- | :--- | :--- |
| **Windows Desktop** | Wails v2/v3 | `led-swarm.exe` (Single Executable with embedded webview) |
| **macOS Desktop** | Wails v2/v3 | `led-swarm.app` / DMG (Universal Binary arm64 + x86_64) |
| **Linux Desktop** | Wails v2/v3 | `led-swarm` (Native GTK3 WebKit app) |
| **Headless Server** | `go build -tags server` | `led-swarm-server` (Zero-dependency binary with embedded Svelte UI) |
| **Docker** | Docker Multi-stage | `ghcr.io/user/led-swarm:latest` (<25MB scratch/alpine container) |

---

*Document compiled and stored under `.docs/architecture.md`.*
