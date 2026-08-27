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
	devMgr.SetDashboardManager(dashboardMgr)

	srv := &Server{
		db:           db,
		hub:          hub,
		wledClient:   wledClient,
		devMgr:       devMgr,
		groupMgr:     groupMgr,
		dashboardMgr: dashboardMgr,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/devices", srv.handleDevices)
	mux.HandleFunc("/api/v1/devices/add", srv.handleAddDevice)
	mux.HandleFunc("/api/v1/devices/name", srv.handleUpdateName)
	mux.HandleFunc("/api/v1/effects", srv.handleEffects)
	mux.HandleFunc("/api/v1/palettes", srv.handlePalettes)

	mux.HandleFunc("/api/v1/groups", srv.handleGroups)
	mux.HandleFunc("/api/v1/scenes", srv.handleScenes)

	mux.HandleFunc("/api/v1/dashboard/items", srv.handleDashboardItems)
	mux.HandleFunc("/api/v1/dashboard/pin", srv.handlePinDashboardItem)
	mux.HandleFunc("/api/v1/dashboard/panel", srv.handleSetDashboardItemPanel)
	mux.HandleFunc("/api/v1/dashboard/panels", srv.handleDashboardPanels)
	mux.HandleFunc("/api/v1/dashboard/panels/", srv.handleDashboardPanels)

	ts := httptest.NewServer(mux)
	t.Cleanup(func() {
		ts.Close()
	})

	return srv, ts
}

func TestServer_RESTEndpoints(t *testing.T) {
	_, ts := createTestServer(t)

	// 1. Test GET /api/v1/devices (Empty list initially)
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

	// 2. Test POST /api/v1/dashboard/panels (Add custom panel)
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

	// 3. Test GET /api/v1/dashboard/panels
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

	// 4. Test POST /api/v1/dashboard/panel (Set item panel ID)
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
}
