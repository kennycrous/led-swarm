# Slice 006 Specification: 2D Visual Layout Canvas Editor

## 1. Overview & Objectives

**Slice 006** introduces an interactive **2D Visual Layout Canvas Editor** in Svelte 5 and Go. Users can drag and position physical WLED light strips on a 2D room grid mapping their real-world physical placement, adjust rotation and scale, view live pixel color mirroring in real time, and trigger 2D spatial coordinate sweep animations.

---

## 2. Technical Design & Subsystem Architecture

### 2.1 Backend Subsystem (Go)
1. **SQLite Database Schema (`backend/src/db.go`)**:
   - Table `canvas_placements`:
     - `device_id` TEXT PRIMARY KEY (FOREIGN KEY references `devices(id)` ON DELETE CASCADE)
     - `pos_x` REAL DEFAULT 100.0
     - `pos_y` REAL DEFAULT 100.0
     - `rotation` REAL DEFAULT 0.0
     - `scale` REAL DEFAULT 1.0
     - `geometry` TEXT DEFAULT 'strip'
     - `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP

2. **Canvas Manager (`backend/src/canvas_manager.go`)**:
   - Manages in-memory thread-safe placement map (`map[string]CanvasPlacement`).
   - Provides methods:
     - `GetPlacements() []CanvasPlacement`
     - `SavePlacement(p CanvasPlacement) error`
     - `BatchSavePlacements(placements []CanvasPlacement) error`
     - `AutoPlaceNewDevice(devID string) CanvasPlacement`
   - Broadcasts `canvas_updated` WebSocket event on changes.

3. **REST & Wails API Endpoints**:
   - `GET /api/v1/canvas/placements` -> Returns all placements
   - `POST /api/v1/canvas/placement` -> Saves single placement
   - `POST /api/v1/canvas/placements/batch` -> Saves batch placements
   - Wails bindings: `GetCanvasPlacements()`, `SaveCanvasPlacement()`, `BatchSaveCanvasPlacements()`.

---

### 2.2 Frontend Subsystem (Svelte 5)
1. **Canvas Store (`frontend/src/lib/stores/canvasStore.svelte.js`)**:
   - Reactive Svelte 5 store managing placements state `$state([])`.
   - Methods to update position (`updatePlacement`), rotation, scale, and persist to server.

2. **Canvas Editor Component (`frontend/src/lib/components/CanvasEditor.svelte`)**:
   - **Cyberpunk Grid Canvas**: Glassmorphic dark slate grid `#06090e` with subtle glowing cyan grid lines, zoom/pan controls, grid snap toggle (10px).
   - **Interactive Strip Elements**:
     - Drag-and-drop position adjustment.
     - Interactive rotation handle or slider (0° - 360°).
     - Scale/length handle.
     - Geometry selector (Strip vs Matrix vs Ring).
   - **Live Pixel Mirroring**:
     - Renders individual LED dots along the strip element using live color data from WLED device state (`$state` / deviceStore).
   - **Spatial Sweep Animation Trigger**:
     - Button triggering a 2D spatial rainbow/cyan sweep animation across all canvas elements.

3. **Navigation Integration**:
   - Add **"2D Room Canvas"** tab with `Grid` / `Map` icon to top header navigation.

---

## 3. TDD & Incremental Implementation Steps

### Phase 1: Test Setup (Red Phase)
- Write Go backend unit tests in `backend/src/canvas_manager_test.go`:
  - `TestDatabase_CanvasPlacements`
  - `TestCanvasManager_PlacementsAndBatchSave`
- Write Vitest frontend tests in `frontend/tests/stores/canvasStore.test.js` & `frontend/tests/components/CanvasEditor.test.js`.

### Phase 2: Backend Implementation (Green Phase)
- Update SQLite schema in `db.go` with `canvas_placements` table.
- Implement `canvas_manager.go`, `server.go` handlers, and `app.go` Wails bindings.
- Run `go test ./src/...` to verify backend tests turn green.

### Phase 3: Frontend Implementation (Green Phase)
- Create `canvasStore.svelte.js` and `CanvasEditor.svelte`.
- Integrate top tab navigation in `App.svelte`.
- Run `npm run test` to verify Vitest tests turn green.

### Phase 4: Verification & Quality Assurance
- Run `make fmt && make lint && make test && make build`.
- Manual verification of drag-and-drop, placement saving, live pixel mirroring, and spatial sweep animation.
