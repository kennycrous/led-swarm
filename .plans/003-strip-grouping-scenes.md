# Implementation Plan 003: Virtual Strip Grouping & Multi-Zone Scenes

## 1. Overview & Objectives

This implementation plan covers **Slice 003: Virtual Strip Grouping & Multi-Zone Scenes**. It enables combining physical WLED light strips into logical virtual groups ("Desk Setup", "Ceiling Cove", "TV Backlight") and capturing multi-device state snapshots as named Scenes ("Cyberpunk Cyan", "Movie Night", "Warm Ambient") for instant 1-click execution across the swarm.

---

## 2. Architecture & Data Flow

```
+-----------------------------------------------------------------------------------+
|                        Svelte 5 Cyberpunk UI / Desktop                            |
|                                                                                   |
|  [ Top Quick Scene Launcher ]  [ Groups & Scenes Tab ]  [ Capture Scene Modal ]  |
|                                         │                                         |
|                   Svelte 5 Group & Scene Store ($state)                           |
+-----------------------------------------+-----------------------------------------+
                                          │
                        WebSocket & REST API Dispatcher
                                          │
                                          v
+-----------------------------------------------------------------------------------+
|                            LED Swarm Go Backend Engine                            |
|                                                                                   |
|  +---------------------+  +----------------------+  +--------------------------+  |
|  | Group & Scene       |  | Concurrent Goroutine |  | SQLite Database          |  |
|  | Manager             |  | Dispatcher           |  | (groups, group_devices,  |  |
|  |                     |  | (sync.WaitGroup)     |  |  scenes tables)          |  |
|  +---------------------+  +----------------------+  +--------------------------+  |
+-----------------------------------------+-----------------------------------------+
                                          │
                           Parallel HTTP / WebSocket Push
                                          v
+-----------------------------------------------------------------------------------+
|                            Target WLED Light Strips                               |
|        [ WLED Strip 1 ]           [ WLED Strip 2 ]          [ WLED Strip N... ]   |
+-----------------------------------------------------------------------------------+
```

---

## 3. Step-by-Step Implementation Tasks

### Task 1: SQLite Schema & Database Handlers (`backend/src/db.go`)
- [ ] Implement database methods for groups:
  - `CreateGroup(name string, description string, deviceIDs []string) (*Group, error)`
  - `UpdateGroup(id string, name string, description string, deviceIDs []string) error`
  - `DeleteGroup(id string) error`
  - `GetGroups() ([]Group, error)`
- [ ] Implement database methods for scenes:
  - `CreateScene(name string, icon string, configJSON string) (*Scene, error)`
  - `DeleteScene(id string) error`
  - `GetScenes() ([]Scene, error)`

### Task 2: Group & Scene Backend Engine (`backend/src/group_manager.go`)
- [ ] Create `backend/src/group_manager.go`:
  - Manage in-memory group and scene caches.
  - Implement `SetGroupState(groupID string, rawState json.RawMessage) error`: Uses `sync.WaitGroup` to dispatch state commands concurrently across all assigned device IPs in parallel goroutines.
  - Implement `CaptureScene(name string, icon string) (*Scene, error)`: Queries current live states of all online devices, serializes state snapshot to JSON, and saves to SQLite.
  - Implement `ApplyScene(sceneID string) error`: Deserializes scene JSON snapshot and dispatches state updates concurrently to target devices.
- [ ] Add REST API endpoints in `backend/src/server.go`:
  - `GET /api/v1/groups`, `POST /api/v1/groups`, `DELETE /api/v1/groups/:id`
  - `POST /api/v1/groups/:id/state` (Batch group state dispatch)
  - `GET /api/v1/scenes`, `POST /api/v1/scenes/capture`, `POST /api/v1/scenes/:id/apply`, `DELETE /api/v1/scenes/:id`
- [ ] Add Wails exposed binding methods in `backend/src/app.go`.

### Task 3: Svelte 5 Group & Scene Store (`frontend/src/lib/stores/groupStore.svelte.js`)
- [ ] Create `frontend/src/lib/stores/groupStore.svelte.js` using Svelte 5 `$state` runes.
- [ ] Reactive store getters: `groups`, `scenes`, `activeGroupId`.
- [ ] Reactive actions:
  - `createGroup(name, description, deviceIds)`
  - `updateGroup(groupId, name, description, deviceIds)`
  - `deleteGroup(groupId)`
  - `setGroupState(groupId, statePayload)`
  - `captureScene(name, icon)`
  - `applyScene(sceneId)`
  - `deleteScene(sceneId)`

### Task 4: Groups & Scenes UI Components (`frontend/src/lib/components/`)
- [ ] Build `GroupCard.svelte`:
  - Displays group name, assigned strip count, and member thumbnails.
  - Group batch controls: Master power toggle, group brightness slider, unified color swatch selector, group WLED effect dropdown.
  - Edit & delete buttons.
- [ ] Build `SceneCard.svelte`:
  - Displays scene icon, name, created date, and 1-click "APPLY SCENE" button.
- [ ] Build `CreateGroupModal.svelte`: Modal with strip selection checkboxes for creating/editing virtual groups.
- [ ] Build `CaptureSceneModal.svelte`: Modal with icon selector (Sparkles, Sun, Moon, Film, Flame, Zap) and scene name input.

### Task 5: App Integration & Top Cyber Bar Launcher (`frontend/src/App.svelte`)
- [ ] Update `App.svelte` header bar to include Top Scene Quick Preset buttons for 1-click scene activation from any tab.
- [ ] Add Dashboard group filter pills allowing filtering main device grid by active group.
- [ ] Render Groups & Scenes tab with grid of active Groups and Scenes.

### Task 6: Roadmap Documentation Update
- [ ] Update `.docs/roadmap.md` to mark Slice 003 as COMPLETED (`🟢 Slice 003 [COMPLETED]`) with checked deliverables.

---

## 4. Verification & Testing Criteria

1. **Virtual Group Creation & Batch Control**:
   - Create group "Desk Setup" containing 2 WLED strips.
   - Sliding group brightness slider changes both strips' brightness simultaneously.
   - Toggling group power turns all assigned strips on/off concurrently.
2. **One-Click Scene Capture & Restoration**:
   - Set Strip 1 to Cyan and Strip 2 to Magenta.
   - Click "Capture Scene" -> Name: "Cyberpunk Night" -> Save.
   - Change both strips to Green.
   - Click "Cyberpunk Night" scene button in top header -> Both strips revert to Cyan and Magenta in <50ms.
3. **Database Persistence**:
   - Restart Go backend server -> Groups, group memberships, and saved scenes remain intact from SQLite.

---

*Plan stored in `.plans/003-strip-grouping-scenes.md`.*
