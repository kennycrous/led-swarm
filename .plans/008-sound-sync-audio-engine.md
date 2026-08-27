# Slice 008 Implementation Plan: Card-Integrated Sound Sync & Audio Reactivity Engine

## Executive Summary
Slice 008 adds real-time audio reactivity and sound synchronization to **LED Swarm Orchestrator**. System audio is captured via WebAudio API in Svelte 5 (or Wails desktop webview) and binned into a 16-band FFT spectrum analyzer. Audio energy vectors (`audio_fft`) are dispatched to the Go backend over WebSocket (`/api/v1/ws`). 

When **Sound Sync Mode** is toggled ON for any **Device Card**, **Group Card**, or **2D Room Canvas Card**, the backend DDP Streamer engine (`backend/src/ddp_streamer.go`) modulates 60 FPS DDP UDP (port 4048) pixel buffers using real-time audio frequency energy (Bass, Mid, Treble) across assigned strips, virtual groups, or 2D room coordinates.

---

## Architectural Decisions

1. **WebAudio API + WebSocket Proxy Engine**:
   - WebAudio API (`navigator.mediaDevices.getUserMedia`) in Svelte UI captures system mic/stereo mix audio and computes 16-band FFT frequency energy vectors.
   - Broadcasts real-time `audio_fft` packets over WebSocket to backend `Hub`.
   - Backend `AudioEngine` caches latest FFT energy (`Bass: 20-250Hz`, `Mid: 250-4000Hz`, `Treble: 4-20kHz`, `Full Spectrum`) and pushes energy updates to active DDP streams.

2. **60 FPS Target-Keyed DDP Sound Reactive Streamer**:
   - Sound Sync mode operates on a target-keyed basis (`"device:<id>"`, `"group:<id>"`, `"room:<id>"`).
   - Generates sound-reactive DDP procedural effects:
     - 🔊 **Audio Bass Pulse**: Low-frequency bass kick pulses center-outward across 1D/2D space.
     - 🌊 **Audio Spectrum Waterfall**: 16-band frequency spectrum mapped linearly across pixels.
     - 🎯 **Audio 2D Beat Ripple**: Radial 2D room canvas ripples triggered on audio beat peaks.
     - 📊 **Audio VU Meter**: Real-time stereo VU level meter bouncing across strip length.
     - ✨ **Audio Treble Sparkle**: High-frequency treble transients trigger bright neon sparkles.

3. **Card-Level Sound Sync UI & Mini FFT Visualizer**:
   - Dedicated `Sound Sync` button (`Volume2` icon) on `DeviceCard.svelte`, `GroupCard.svelte`, and `RoomCard.svelte` headers next to `DDP 60`.
   - Toggling Sound Sync opens an accordion drawer featuring:
     - Mini live 16-band FFT spectrum visualizer bar (`cyan`/`purple`/`gold` bars).
     - Audio Sensitivity Gain Slider ($0.1\times - 5.0\times$).
     - Frequency Band Selector (`Bass`, `Mid`, `Treble`, `Full Spectrum`).
     - Sound-Reactive Effect Preset Dropdown.

---

## Implementation Steps & Tasks

### Step 1: Frontend Audio Capture & FFT Store (`frontend/src/lib/stores/audioStore.svelte.js`)
- [ ] Create `audioStore.svelte.js` managing WebAudio `AudioContext`, `AnalyserNode`, and `MediaStream`.
- [ ] Compute 16-band frequency bin energy (normalized 0.0 - 1.0) at 60 FPS using `requestAnimationFrame`.
- [ ] Send `audio_fft` WebSocket messages to backend `Hub` containing 16-band energy array.
- [ ] Unit test `audioStore.test.js`.

### Step 2: Backend Audio Engine & WS Dispatcher (`backend/src/audio_engine.go`)
- [ ] Implement `AudioEngine` in Go caching latest 16-band FFT energy frame.
- [ ] Calculate frequency band energy helpers (`GetBass()`, `GetMid()`, `GetTreble()`, `GetPeak()`).
- [ ] Update WebSocket `Hub` to handle `audio_fft` client messages and forward energy state to `ddpStreamer`.
- [ ] Unit test `audio_engine_test.go`.

### Step 3: Sound-Reactive DDP Procedural Generators (`backend/src/ddp_streamer.go`)
- [ ] Update `DDPStreamer` to accept audio energy parameters per active target stream.
- [ ] Implement 1D & 2D sound-reactive procedural effect algorithms:
  - `audio_bass_pulse`
  - `audio_spectrum_waterfall`
  - `audio_spatial_ripple`
  - `audio_vu_meter`
  - `audio_treble_sparkle`
- [ ] Update `ddp_streamer_test.go`.

### Step 4: Card UI Integration (`DeviceCard.svelte`, `GroupCard.svelte`, `RoomCard.svelte`)
- [ ] Add `Sound Sync` toggle button (`Volume2`) to `DeviceCard.svelte`, `GroupCard.svelte`, and `RoomCard.svelte`.
- [ ] Add expandable Sound Sync drawer with mini FFT bar chart, gain slider, frequency band selector, and audio preset dropdown.
- [ ] Pass `audioStore` to all cards in `App.svelte`.
- [ ] Component unit tests (`DeviceCard.test.js`, `GroupCard.test.js`, `RoomCard.test.js`).

### Step 5: Verification & Quality Assurance
- [ ] Execute `make fmt && make lint && make test && make build`.
- [ ] Verify 100% test pass rate across backend Go tests and frontend Vitest suite.

---

*Plan specification saved under `.plans/008-sound-sync-audio-engine.md`.*
