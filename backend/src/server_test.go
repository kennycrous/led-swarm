package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestServer(t *testing.T) (*Server, *httptest.Server) {
	t.Helper()

	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()

	devMgr := NewDeviceManager(db, wledClient, hub)
	groupMgr := NewGroupManager(db, wledClient, devMgr, hub)
	dashboardMgr := NewDashboardManager(db, hub)
	canvasMgr := NewCanvasManager(db, hub)
	ddpStreamer := NewDDPStreamer()
	ddpStreamer.SetHub(hub)
	devMgr.SetDashboardManager(dashboardMgr)

	srv := &Server{
		db:           db,
		hub:          hub,
		wledClient:   wledClient,
		devMgr:       devMgr,
		groupMgr:     groupMgr,
		dashboardMgr: dashboardMgr,
		canvasMgr:    canvasMgr,
		ddpStreamer:  ddpStreamer,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/devices", srv.handleDevices)
	mux.HandleFunc("/api/v1/devices/add", srv.handleAddDevice)
	mux.HandleFunc("/api/v1/devices/name", srv.handleUpdateName)
	mux.HandleFunc("/api/v1/effects", srv.handleEffects)
	mux.HandleFunc("/api/v1/palettes", srv.handlePalettes)

	mux.HandleFunc("/api/v1/groups", srv.handleGroups)
	mux.HandleFunc("/api/v1/groups/", srv.handleGroups)
	mux.HandleFunc("/api/v1/scenes", srv.handleScenes)

	mux.HandleFunc("/api/v1/dashboard/items", srv.handleDashboardItems)
	mux.HandleFunc("/api/v1/dashboard/pin", srv.handlePinDashboardItem)
	mux.HandleFunc("/api/v1/dashboard/panel", srv.handleSetDashboardItemPanel)
	mux.HandleFunc("/api/v1/dashboard/panels", srv.handleDashboardPanels)
	mux.HandleFunc("/api/v1/dashboard/panels/", srv.handleDashboardPanels)

	mux.HandleFunc("/api/v1/canvas/rooms", srv.handleCanvasRooms)
	mux.HandleFunc("/api/v1/canvas/rooms/", srv.handleCanvasRooms)

	ts := httptest.NewServer(mux)
	t.Cleanup(func() {
		ts.Close()
	})

	return srv, ts
}

func TestServer_RESTEndpoints(t *testing.T) {
	_, ts := createTestServer(t)

	// 1. GET /api/v1/devices
	res, err := http.Get(ts.URL + "/api/v1/devices")
	if err != nil {
		t.Fatalf("GET /api/v1/devices failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var devices []Device
	json.NewDecoder(res.Body).Decode(&devices)
	res.Body.Close()
	if len(devices) != 0 {
		t.Errorf("Expected 0 devices initially, got %d", len(devices))
	}

	// 2. GET /api/v1/effects & /api/v1/palettes
	res, _ = http.Get(ts.URL + "/api/v1/effects")
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for /api/v1/effects, got %d", res.StatusCode)
	}
	res.Body.Close()

	res, _ = http.Get(ts.URL + "/api/v1/palettes")
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for /api/v1/palettes, got %d", res.StatusCode)
	}
	res.Body.Close()

	// 3. POST /api/v1/dashboard/panels
	panelReq := map[string]string{"title": "Bedroom Ambient"}
	body, _ := json.Marshal(panelReq)
	res, err = http.Post(ts.URL+"/api/v1/dashboard/panels", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("POST /api/v1/dashboard/panels failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var createdPanel DashboardPanel
	json.NewDecoder(res.Body).Decode(&createdPanel)
	res.Body.Close()
	if createdPanel.Title != "Bedroom Ambient" {
		t.Errorf("Expected title 'Bedroom Ambient', got '%s'", createdPanel.Title)
	}

	// 4. GET /api/v1/dashboard/panels
	res, err = http.Get(ts.URL + "/api/v1/dashboard/panels")
	if err != nil {
		t.Fatalf("GET /api/v1/dashboard/panels failed: %v", err)
	}
	var panels []DashboardPanel
	json.NewDecoder(res.Body).Decode(&panels)
	res.Body.Close()
	if len(panels) != 1 {
		t.Fatalf("Expected 1 panel, got %d", len(panels))
	}

	// 5. POST /api/v1/dashboard/panel
	itemReq := map[string]string{"itemId": "wled-test-99", "panelId": createdPanel.ID}
	body, _ = json.Marshal(itemReq)
	res, err = http.Post(ts.URL+"/api/v1/dashboard/panel", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("POST /api/v1/dashboard/panel failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var item DashboardItem
	json.NewDecoder(res.Body).Decode(&item)
	res.Body.Close()
	if item.PanelID != createdPanel.ID {
		t.Errorf("Expected PanelID '%s', got '%s'", createdPanel.ID, item.PanelID)
	}

	// 6. Test Canvas Room Creation & Renaming REST endpoints
	roomReq := map[string]interface{}{"title": "Living Room", "width": 2000, "height": 1200}
	body, _ = json.Marshal(roomReq)
	res, err = http.Post(ts.URL+"/api/v1/canvas/rooms", "application/json", bytes.NewBuffer(body))
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("POST /api/v1/canvas/rooms failed: %v, status: %d", err, res.StatusCode)
	}
	var createdRoom CanvasRoom
	json.NewDecoder(res.Body).Decode(&createdRoom)
	res.Body.Close()

	// Now rename the room via POST /api/v1/canvas/rooms/rename
	renameReq := map[string]string{"id": createdRoom.ID, "title": "Main Living Area"}
	body, _ = json.Marshal(renameReq)
	res, err = http.Post(ts.URL+"/api/v1/canvas/rooms/rename", "application/json", bytes.NewBuffer(body))
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("POST /api/v1/canvas/rooms/rename failed: %v, status: %d", err, res.StatusCode)
	}
	var renamedRoom CanvasRoom
	json.NewDecoder(res.Body).Decode(&renamedRoom)
	res.Body.Close()

	if renamedRoom.ID != createdRoom.ID {
		t.Errorf("Expected same Room ID '%s', got '%s'", createdRoom.ID, renamedRoom.ID)
	}
	if renamedRoom.Title != "Main Living Area" {
		t.Errorf("Expected title 'Main Living Area', got '%s'", renamedRoom.Title)
	}

	// 6. Test Error Cases (HTTP 400 Bad Request)
	res, _ = http.Post(ts.URL+"/api/v1/devices/add", "application/json", bytes.NewBufferString(`{"ip":""}`))
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for empty IP, got %d", res.StatusCode)
	}
	res.Body.Close()

	res, _ = http.Post(ts.URL+"/api/v1/dashboard/panels", "application/json", bytes.NewBufferString(`{"title":""}`))
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for empty panel title, got %d", res.StatusCode)
	}
	res.Body.Close()
}
