# Implementation Plan 002: Live WLED Device Discovery & Real-Time Controls

## 1. Overview & Objectives

This implementation plan covers **Slice 002: Live WLED Device Discovery & Real-Time Controls**. It connects the Svelte 5 frontend to live physical WLED light strips on the local network, replacing sample UI data with bidirectional WebSocket state push, real-time controls (power, brightness, RGB color, WLED effects & palettes), custom device nicknames, and SQLite persistence.

---

## 2. Technical Architecture & Data Flow

```
+-----------------------------------------------------------------------------------+
|                        Svelte 5 Cyberpunk Web UI / Desktop                        |
|                                                                                   |
|  [ Master Swarm Bar ]   [ Device Card ]   [ Color Picker ]   [ Effect Selector ]  |
|                                     │                                             |
|                   Unified Svelte 5 State Store ($state)                           |
+-------------------------------------+---------------------------------------------+
                                      │
                     WebSocket / WS API & REST Fallback
                                      │
                                      v
+-----------------------------------------------------------------------------------+
|                            LED Swarm Go Backend Engine                            |
|                                                                                   |
|  +---------------------+  +----------------------+  +--------------------------+  |
|  |  Device Manager     |  |  WLED WebSocket Pool |  | SQLite Database          |  |
|  |  State Cache        |  |  (ws://<ip>/ws)      |  | (devices table)          |  |
|  +---------------------+  +----------------------+  +--------------------------+  |
|  | Effect/Palette Cache|  |  WebSocket Hub       |  | mDNS Scanner             |  |
|  | (/api/v1/effects)   |  |  (/api/v1/ws)        |  | (_wled._tcp)             |  |
|  +---------------------+  +----------------------+  +--------------------------+  |
+-------------------------------------+---------------------------------------------+
                                      │
                      HTTP / WS WLED JSON API Protocol
                                      v
+-----------------------------------------------------------------------------------+
|                            Target WLED Light Strips                               |
|        [ WLED Strip 1 ]           [ WLED Strip 2 ]          [ WLED Strip N... ]   |
+-----------------------------------------------------------------------------------+
```

---

## 3. Step-by-Step Implementation Tasks

### Task 1: Go WLED Device Manager & WebSocket Client Pool (`backend/src/`)
- [ ] Create WLED device manager module (`backend/src/device_manager.go`) to manage live connection state for all discovered WLED devices.
- [ ] Implement persistent WLED WebSocket connections (`ws://<ip>/ws`) in Go with automatic reconnection on disconnect.
- [ ] Implement state update dispatcher: When WLED device pushes state update over WebSocket, update local Go cache and broadcast delta to Svelte UI clients via WebSocket Hub (`hub.go`).
- [ ] Add REST API endpoints:
  - `POST /api/v1/devices/add` (Manually add WLED device by IP address)
  - `POST /api/v1/devices/:id/name` (Update device custom nickname in SQLite)
  - `GET /api/v1/effects` (Returns cached WLED effect names list)
  - `GET /api/v1/palettes` (Returns cached WLED palette names list)

### Task 2: Database Layer Enhancements (`backend/src/db.go`)
- [ ] Ensure SQLite `devices` table stores `custom_name`, `ip_address`, `mac_address`, `led_count`, `is_online`, `last_seen`, and last known `state_json`.
- [ ] Add database methods:
  - `UpdateDeviceName(id string, customName string) error`
  - `DeleteDevice(id string) error`
  - `GetDeviceByID(id string) (*Device, error)`

### Task 3: Unified Svelte 5 State Store (`frontend/src/lib/stores/`)
- [ ] Create `frontend/src/lib/stores/deviceStore.svelte.js` using Svelte 5 `$state` runes.
- [ ] Implement environment detection:
  - If running in Wails desktop app (`window.runtime` present), use Wails Go bindings (`GetDevices()`, `SetDevicePower()`, `SetDeviceBrightness()`).
  - If running in Web browser / Docker mode, open WebSocket to `ws://${location.host}/api/v1/ws` and fetch `/api/v1/devices`.
- [ ] Provide reactive state methods:
  - `togglePower(deviceId)`
  - `setBrightness(deviceId, value)`
  - `setColor(deviceId, rgb)`
  - `setEffect(deviceId, effectId)`
  - `setPalette(deviceId, paletteId)`
  - `setMasterPower(on)`
  - `setMasterBrightness(value)`
  - `manualAddDevice(ip)`

### Task 4: Cyberpunk UI Components (`frontend/src/lib/components/`)
- [ ] Build `DeviceCard.svelte`:
  - Glowing status indicator pill (Emerald = Online, Crimson = Offline).
  - Power toggle button with neon cyan glow.
  - Per-device brightness slider with percentage readout.
  - Color spectrum picker popover / preset color swatches.
  - WLED Effect & Palette dropdown select inputs.
  - Nickname edit modal / inline edit.
- [ ] Build `ManualIpModal.svelte`: Modal for manually entering WLED device IP address.
- [ ] Update `App.svelte` to bind `deviceStore` reactively, replacing sample data.

### Task 5: Wails Desktop Bindings Integration (`backend/src/app.go`)
- [ ] Expose new Wails methods on `App` struct:
  - `SetDeviceColor(ip string, r int, g int, b int) error`
  - `SetDeviceEffect(ip string, fx int) error`
  - `SetDevicePalette(ip string, pal int) error`
  - `AddManualDevice(ip string) (*Device, error)`
  - `UpdateDeviceNickname(id string, name string) error`

---

## 4. Verification & Testing Steps

1. **Local Network WLED Discovery**:
   - Run `make dev-backend` on a network with WLED light strips.
   - Verify mDNS auto-discovers WLED devices and populates SQLite.
2. **Real-Time Control Testing**:
   - Toggling power in Svelte UI instantly turns the physical WLED strip on/off.
   - Adjusting brightness slider changes physical LED intensity smoothly.
   - Picking an RGB color updates physical WLED color without delay.
3. **Bidirectional State Push**:
   - Toggling power via the WLED physical button or native WLED mobile app instantly updates the Svelte UI state via WebSocket.
4. **Manual IP Addition**:
   - Manually entering a WLED IP in the modal connects to the strip and saves it in SQLite.
5. **Wails Desktop & Web Mode**:
   - Verify functionality works cleanly in both `wails dev` desktop mode and `make dev-backend` server mode.

---

*Plan stored in `.plans/002-device-discovery-controls.md`.*
