# AGENTS.md - Repository Overview & AI Agent Guidelines

Welcome to **LED Swarm Orchestrator** (`led-swarm`). This document provides AI coding agents with a comprehensive overview of the repository, architectural decisions, technical stack, and design guidelines established for this project.

---

## 1. Project Overview & Objectives

**LED Swarm** is a modern, high-performance WLED orchestrator application and server designed to control, manage, and synchronize multiple WLED light strips.

### Core Problems Solved:
- Outdated, unmaintained, or slow existing WLED management software.
- High memory usage and complex dependencies in Electron-based LED apps.
- Lack of unified multi-strip synchronization, 2D room canvas layout mapping, low-latency DDP pixel streaming, and sound reactivity.

### Dual Operating Modes:
1. **Native Desktop Application**: Self-contained OS desktop app using Wails v2/v3 (~15MB executable size, low RAM usage).
2. **Headless Server & Docker**: Single zero-dependency Go binary with embedded web UI served via HTTP (`--server` mode or containerized in Docker).

---

## 2. Technical Stack

| Component | Technology / Library | Rationale |
| :--- | :--- | :--- |
| **Backend Engine** | Go 1.22+ | Fast concurrency, native networking, low memory footprint, zero-dependency static binaries. |
| **Database** | CGO-free Pure-Go SQLite (`modernc.org/sqlite` / `glebarez/go-sqlite`) | Zero C toolchain required for cross-compilation across Windows, macOS, Linux, and ARM/Raspberry Pi. |
| **Frontend Framework** | Svelte 5 + Vite | Ultra-lightweight (~30KB bundle), zero virtual DOM overhead, high reactivity for live canvas and controls. |
| **Styling & Theme** | Tailwind CSS + Lucide Icons | Dark Cyberpunk Glassmorphism (`backdrop-blur-xl`, void dark background `#06090e`, cyan/purple neon accents). |
| **Desktop Packaging** | Wails v2/v3 | Go native desktop framework using system webview (WebView2/WebKit). |
| **Server Packaging** | Go `go:embed` | Compiles static frontend assets directly into the Go executable for single-binary distribution or Docker (<25MB). |
| **WLED Protocols** | mDNS (`_wled._tcp`), JSON REST/WebSocket API, DDP (UDP 4048), WLED Realtime UDP (21324) | Auto-discovery, bidirectional state sync, 60 FPS raw pixel streaming. |

---

## 3. Directory Structure & Key Files

```
led-swarm/
├── .docs/
│   ├── architecture.md       # Complete architectural & system design specification
│   └── roadmap.md            # Incremental horizontal feature roadmap specification
├── .plans/                   # Step-by-step implementation plans & phase specifications
├── AGENTS.md                  # AI agent context & guidelines (this file)
└── README.md                  # Project intro
```

---

## 4. Key Subsystems Overview

1. **Auto-Discovery (`mDNS`)**: Background worker scanning for `_wled._tcp` services, auto-registering IPs, MAC addresses, and LED counts.
2. **State Orchestration & WS Hub**: Central Go state engine maintaining device state caches and broadcasting real-time updates to UI clients over WebSockets.
3. **2D Visual Layout Canvas**: Interactive 2D workspace in Svelte for dragging & positioning LED strips visually, featuring real-time pixel mirroring and spatial effects.
4. **Virtual Grouping & Multi-Zone Scenes**: Combine multiple strips into groups and trigger synchronized scenes, brightness, colors, or palettes.
5. **DDP / UDP Realtime Pixel Streamer**: High-frequency 60 FPS pixel buffer generator for custom animations and matrix visuals.
6. **Sound Sync Engine**: Audio capture + FFT binning to map audio energy to LED spatial animations.

---

## 5. UI/UX Design Guidelines (Cyberpunk Glassmorphism)

Agents modifying or adding UI components must adhere to the following design system:

- **Colors**:
  - Background: `#06090e` (Void Black/Dark Slate)
  - Card Surfaces: `rgba(15, 23, 42, 0.65)` with `backdrop-filter: blur(16px)`
  - Borders: Thin glowing borders `rgba(56, 189, 248, 0.2)` with purple glow on hover (`rgba(168, 85, 247, 0.4)`)
  - Accent Colors: Cyan (`#06b6d4`), Magenta/Purple (`#a855f7`), Neon Gold (`#f59e0b`), Online Emerald (`#10b981`)
- **Interactions**: Fluid transitions, hover glow effects, micro-animations, glassmorphic card overlays, responsive canvas elements.

---

## 6. AI Agent Development Guidelines

- **Architecture & Implementation Plans**: Before creating new subsystems or API endpoints, refer to [.docs/architecture.md](file:///mnt/Data/source/led-swarm/.docs/architecture.md) and check `.plans/` for step-by-step implementation plans.
- **CGO Avoidance**: Maintain CGO-free pure Go builds to ensure cross-platform single binary compilation without external C compilers.
- **Embedded Web Assets**: Ensure all Svelte frontend build outputs in `frontend/dist` can be cleanly embedded via `go:embed`.
- **Wails Bindings**: Place Go application methods exposed to the desktop UI in the Wails app struct for auto-generated TypeScript bindings.
- **Test-Driven Development (TDD)**: When implementing any code changes or new features, first write a failing unit or component test (Red phase), implement the code changes, and verify that the test turns green (`make test`) before proceeding.
- **Pre-Commit Pause Directive**: ALWAYS pause and wait for explicit user confirmation before committing any changes (`git commit`), allowing the user to locally test the working tree first.
- **File Links**: Always link to relevant files using standard markdown links (e.g. `[architecture.md](file:///mnt/Data/source/led-swarm/.docs/architecture.md)`).
