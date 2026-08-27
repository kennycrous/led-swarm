# 🌌 LED Swarm Orchestrator

> **High-Performance WLED Swarm Manager & Realtime Pixel Orchestrator**
> 
> *Control, map, synchronize, and stream high-framerate visuals across multiple WLED light strips through a dark, cyberpunk glassmorphic UI or headless Docker server.*

---

## ⚡ Features & Highlights

- 🔍 **Zero-Config Auto-Discovery**: Background mDNS worker (`_wled._tcp`) automatically detects WLED light strips on your local network, mapping IP addresses, MAC addresses, and LED counts.
- ⚡ **Low-Latency Pixel Streaming**: Streams 60 FPS raw RGB frame buffers over **DDP (Distributed Display Protocol UDP port 4048)** and WLED Realtime UDP protocol without EEPROM wear.
- 🎨 **Cyberpunk Glassmorphic UI**: Svelte 5 + Tailwind CSS web interface featuring dark slate surfaces (`#06090e`), neon cyan/magenta glows, fluid micro-interactions, and real-time state mirroring.
- 🗺️ **2D Visual Layout Canvas**: Interactive 2D workspace to position WLED strips according to physical room layout for spatial sweep effects.
- 🎛️ **Virtual Groups & Multi-Zone Scenes**: Group multiple strips together into synced zones ("Desk", "Ceiling", "TV") and trigger synchronized scenes with one click.
- 🎵 **Sound Sync & Audio Reactivity**: Integrated FFT frequency binning mapping audio energy dynamically to LED brightness and color pulses.
- 📦 **Zero-Dependency Single Binary**: Pure-Go SQLite engine (`modernc.org/sqlite` / `glebarez/go-sqlite`) requiring zero C toolchains for cross-platform desktop executables or tiny Docker containers (<25MB).

---

## 🚀 Dual Operating Modes

1. **Native Desktop Application**: Self-contained cross-platform OS app built with **Wails v2/v3** using the native system webview (~15MB binary, low RAM usage).
2. **Headless Server & Docker Mode**: Single static Go binary embedding the Svelte 5 web UI (`go:embed`). Runs directly on Linux/macOS/Windows or containerized in Docker / Raspberry Pi.

---

## 🛠️ Tech Stack

| Subsystem | Technology | Description |
| :--- | :--- | :--- |
| **Backend Engine** | Go 1.22+ | High-concurrency mDNS scanner, WebSocket hub, and DDP streamer |
| **Database** | Pure-Go SQLite (`glebarez/go-sqlite`) | CGO-free zero-config embedded database |
| **Frontend UI** | Svelte 5 + Vite + Tailwind CSS | Ultra-lightweight (~30KB bundle), zero Virtual DOM overhead |
| **Icons** | Lucide Icons (`lucide-svelte`) | Cyberpunk aesthetic vectors |
| **Desktop Packaging** | Wails v2/v3 | Native OS WebView container (WebKit / WebView2) |
| **Server Packaging** | Go `go:embed` & Docker | Multi-stage Docker build producing <25MB runner image |

---

## 💻 Quick Start & Local Development

Prerequisites: **Node.js 20+** and **Go 1.22+**.

### Development Commands (`Makefile`)

```bash
# Clone the repository
git clone https://github.com/kennycrous/led-swarm.git
cd led-swarm

# 1. Run Svelte 5 Frontend UI in dev mode with Hot Module Replacement (http://localhost:5173)
make dev-frontend

# 2. Run Go Backend Server locally in dev mode (http://localhost:8080)
make dev-backend

# 3. Run Native Wails Desktop application with live reload
make dev-desktop
```

---

## 📦 Building for Production

### Standalone Executable Build
```bash
# Builds frontend assets into backend/src/dist and compiles the static Go server binary
make build

# Run the compiled headless server binary
./backend/led-swarm-server --server --port=8080
```

### Docker Container Deployment
```bash
# Build Docker image (<25MB)
docker build -t led-swarm:latest .

# Run container on port 8080
docker run -d -p 8080:8080 --name led-swarm led-swarm:latest
```

---

## 📚 Documentation & Architecture Specs

- 📑 [System Architecture Specification](file:///.docs/architecture.md) - Deep-dive into backend engine design, DDP streaming protocol, database schema, and UI tokens.
- 🤖 [AI Agent Context & Guidelines](file:///AGENTS.md) - Context and design principles for AI pair programming.
- 📋 [Implementation Plans](file:///.plans/001-initialization.md) - Step-by-step development roadmap.

---

## 📄 License

MIT License. Developed for WLED lighting enthusiasts.
