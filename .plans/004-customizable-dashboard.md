# Implementation Plan 004: Customizable Dashboard & Device Visibility Management

## 1. Overview & Objectives

This implementation plan covers **Slice 004: Customizable Dashboard & Device Visibility Management**. It allows users to personalize their main Cyber Dashboard layout by pinning/unpinning individual WLED light strips, virtual groups, or scene presets. It also creates a dedicated **Strips & Devices Management** view inside Settings where all network WLED light strips are listed regardless of whether they are visible/pinned on the main dashboard.

---

## 2. Architecture & Data Flow

```
+-----------------------------------------------------------------------------------+
|                        Svelte 5 Cyberpunk UI / Desktop                            |
|                                                                                   |
|  [ Customizable Dashboard Grid ]  [ Settings -> Strips & Devices Management ]    |
|   - Pinned Strips, Groups &       - Full network WLED inventory list             |
|     Scene Presets                  - Pin / Unpin toggles, Rename, Add, Delete     |
|                                         │                                         |
|                 Svelte 5 Dashboard Store ($state runes)                           |
+-----------------------------------------+-----------------------------------------+
                                          │
                        WebSocket & REST API Dispatcher
                                          │
                                          v
+-----------------------------------------------------------------------------------+
|                            LED Swarm Go Backend Engine                            |
|                                                                                   |
|  +---------------------+  +----------------------+  +--------------------------+  |
|  | Dashboard Manager   |  | Auto-Pinning Engine  |  | SQLite Database          |  |
|  |                     |  | (On mDNS discovery)  |  | (dashboard_items table)  |  |
|  +---------------------+  +----------------------+  +--------------------------+  |
+-----------------------------------------------------------------------------------+
```

---

## 3. Step-by-Step Implementation Tasks

### Task 1: SQLite Schema & Database Handlers (`backend/src/db.go`)
- [ ] Add `dashboard_items` SQLite table schema:
  - `id` (TEXT PRIMARY KEY), `item_id` (TEXT NOT NULL), `item_type` (TEXT NOT NULL: 'device'|'group'|'scene'), `position` (INTEGER DEFAULT 0), `is_pinned` (BOOLEAN DEFAULT TRUE), `created_at` (DATETIME).
- [ ] Add database methods:
  - `PinDashboardItem(itemID string, itemType string, isPinned bool) error`
  - `GetDashboardItems() ([]DashboardItem, error)`
  - `ReorderDashboardItems(items []DashboardItem) error`

### Task 2: Dashboard Backend Manager (`backend/src/dashboard_manager.go`)
- [ ] Create `backend/src/dashboard_manager.go`:
  - Manage in-memory pinned items cache.
  - Auto-pinning hook: When a device or group is created/discovered, auto-create a pinned `dashboard_items` entry.
  - REST API routes in `backend/src/server.go`:
    - `GET /api/v1/dashboard/items`
    - `POST /api/v1/dashboard/pin`
    - `POST /api/v1/dashboard/reorder`
  - Wails desktop binding methods in `backend/src/app.go`.

### Task 3: Svelte 5 Dashboard Store (`frontend/src/lib/stores/dashboardStore.svelte.js`)
- [ ] Create `frontend/src/lib/stores/dashboardStore.svelte.js` using `$state` runes.
- [ ] Actions: `init()`, `togglePin(itemId, itemType)`, `reorder(items)`.

### Task 4: Card Pin/Unpin Controls (`frontend/src/lib/components/`)
- [ ] Add Pin/Unpin pin icon toggle to `DeviceCard.svelte`, `GroupCard.svelte`, `SceneCard.svelte`.

### Task 5: Settings -> Strips & Devices Management View (`frontend/src/lib/components/StripsManagement.svelte`)
- [ ] Build `StripsManagement.svelte`:
  - Complete inventory table/cards of ALL discovered and saved WLED strips on the network.
  - Displays IP, MAC address, LED count, online status pill, and inline nickname editor.
  - Per-strip controls: Pin/unpin from dashboard toggle, manual IP adder button, mDNS scan button, and delete/forget device button.

### Task 6: Main Dashboard Custom Grid (`frontend/src/App.svelte`)
- [ ] Update `App.svelte` main Dashboard to render a unified grid of pinned items (Mix of Individual Strip Cards, Group Cards, and 1-Click Scene Cards).
- [ ] Integrate Settings tab with `StripsManagement.svelte`.

### Task 7: Roadmap Documentation Update
- [ ] Update `.docs/roadmap.md` to mark Slice 004 as COMPLETED (`🟢 Slice 004 [COMPLETED]`) upon successful verification.

---

## 4. Verification & Testing Criteria

1. **Dashboard Customization & Pinning**:
   - Create a group "Living Room" containing 2 strips.
   - Unpin the 2 individual strips from the Dashboard -> Dashboard shows only the "Living Room" group card.
2. **Settings -> Strips & Devices Inventory**:
   - Navigate to Settings -> Strips & Devices Management.
   - All network strips remain visible in Settings regardless of whether they are unpinned/hidden from the main Dashboard.
3. **Auto-Pinning on Discovery**:
   - Discover or manually add a new WLED IP -> New strip automatically pins to the main Dashboard by default.

---

*Plan stored in `.plans/004-customizable-dashboard.md`.*
