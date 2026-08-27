# Slice 006 Expanded Specification: Multi-Room 2D Layout Canvas Editor, Scoped Scenes & Room Dashboard Cards

## 1. Overview & Objectives

**Slice 006** introduces an interactive **Multi-Room 2D Visual Layout Canvas Editor**, **Room-Scoped Scenes**, and **Pinnable 2D Room Cards** for the Dashboard. Users can:
1. Create and manage multiple 2D Room Canvases (e.g. *Living Room*, *Office*, *Bedroom*).
2. Position physical WLED light strips on room grid coordinates $(x, y, \text{rotation}, \text{scale}, \text{geometry})$.
3. Pin 2D Room Canvas cards directly to the main Dashboard grid alongside individual strips, groups, and scenes.
4. Capture & apply **Scoped Scenes** targeted at specific Rooms or Groups so that activating a scene in one room does not override other rooms.
5. Create different types of collection groups (*Standard Group*, *2D Spatial Room*, *Scoped Scene*).

---

## 2. Technical Design & Subsystem Architecture

### 2.1 Backend Subsystem (Go)
1. **SQLite Database Schema Extensions (`backend/src/db.go`)**:
   - Table `canvas_rooms`:
     - `id` TEXT PRIMARY KEY
     - `title` TEXT NOT NULL
     - `width` INTEGER DEFAULT 2000
     - `height` INTEGER DEFAULT 1200
     - `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
   - Table `canvas_placements`:
     - `device_id` TEXT NOT NULL
     - `room_id` TEXT NOT NULL DEFAULT 'default'
     - `pos_x` REAL DEFAULT 100.0
     - `pos_y` REAL DEFAULT 100.0
     - `rotation` REAL DEFAULT 0.0
     - `scale` REAL DEFAULT 1.0
     - `geometry` TEXT DEFAULT 'strip'
     - PRIMARY KEY (`device_id`, `room_id`)
   - Table `scenes` updates:
     - `scope_type` TEXT DEFAULT 'global' ('global', 'room', 'group')
     - `target_id` TEXT DEFAULT ''

2. **Canvas Manager (`backend/src/canvas_manager.go`)**:
   - Manages rooms and placement maps.
   - Provides methods for rooms CRUD and per-room placements.

3. **Dashboard & Group Manager Extensions (`dashboard_manager.go`, `group_manager.go`)**:
   - Supports `item_type: 'room'` on Dashboard.
   - Group Creation supporting types (`standard`, `spatial_room`, `scoped_scene`).
   - Scoped scene capture & restoration filter.

---

### 2.2 Frontend Subsystem (Svelte 5)
1. **Room Canvas Selector**:
   - Create, select, rename, and delete 2D Room Canvases in `CanvasEditor.svelte`.
2. **Pinnable Room Card (`RoomCard.svelte`)**:
   - Rendered on the Dashboard grid when pinned.
   - Features mini 2D room map preview, glowing LED dots, master power/brightness, and "Edit 2D Layout" button.
3. **Group Creation & Type Selector Modal**:
   - Choose between *Standard Group*, *2D Spatial Room*, and *Scoped Scene*.
4. **Scoped Scene Capture Modal**:
   - Select scope (*Global Swarm*, *Specific Room*, or *Specific Group*).

---

## 3. TDD & Implementation Phases

### Phase 1: Test Setup (Red Phase)
- Write Go unit tests in `backend/src/canvas_manager_test.go` and `backend/src/group_manager_test.go`:
  - `TestCanvasManager_RoomsAndPlacementsCRUD`
  - `TestGroupManager_ScopedSceneRestore`
- Write Vitest frontend tests in `frontend/tests/components/RoomCard.test.js` & `frontend/tests/stores/canvasStore.test.js`.

### Phase 2: Backend Implementation (Green Phase)
- Implement schema updates, `canvas_rooms` CRUD, scoped scene filtering.
- Verify `go test ./src/...` turns green.

### Phase 3: Frontend Implementation (Green Phase)
- Create `RoomCard.svelte`, update `CanvasEditor.svelte` with room selector, update `App.svelte` dashboard rendering.
- Verify `npm run test` turns green.

### Phase 4: Verification
- Run `make fmt && make lint && make test && make build`.
