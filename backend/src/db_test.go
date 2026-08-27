package main

import (
	"path/filepath"
	"testing"
)

func setupTestDB(t *testing.T) *Database {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestDatabase_Devices(t *testing.T) {
	db := setupTestDB(t)

	// 1. Test SaveDevice & GetDevices
	dev := Device{
		ID:         "wled-test-1",
		Name:       "Test Strip 1",
		IPAddress:  "192.168.1.100",
		MACAddress: "AA:BB:CC:DD:EE:FF",
		LEDCount:   60,
		IsOnline:   true,
	}

	if err := db.SaveDevice(dev); err != nil {
		t.Fatalf("SaveDevice failed: %v", err)
	}

	devices, err := db.GetDevices()
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}
	if len(devices) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(devices))
	}
	if devices[0].Name != "Test Strip 1" {
		t.Errorf("Expected device name 'Test Strip 1', got '%s'", devices[0].Name)
	}

	// 2. Test UpdateDeviceName
	if err := db.UpdateDeviceName("wled-test-1", "Renamed Strip"); err != nil {
		t.Fatalf("UpdateDeviceName failed: %v", err)
	}

	devices, _ = db.GetDevices()
	if devices[0].Name != "Renamed Strip" {
		t.Errorf("Expected renamed device 'Renamed Strip', got '%s'", devices[0].Name)
	}
}

func TestDatabase_GroupsAndScenes(t *testing.T) {
	db := setupTestDB(t)

	// 1. Test SaveGroup & GetGroups
	group := Group{
		ID:          "group-1",
		Name:        "Desk Setup",
		Description: "All desk strips",
		DeviceIDs:   []string{"wled-1", "wled-2"},
	}

	if err := db.SaveGroup(group); err != nil {
		t.Fatalf("SaveGroup failed: %v", err)
	}

	groups, err := db.GetGroups()
	if err != nil {
		t.Fatalf("GetGroups failed: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(groups))
	}
	if len(groups[0].DeviceIDs) != 2 {
		t.Errorf("Expected 2 device IDs in group, got %d", len(groups[0].DeviceIDs))
	}

	// 2. Test SaveScene & GetScenes
	scene := Scene{
		ID:         "scene-1",
		Name:       "Cyberpunk Neon",
		Icon:       "Sparkles",
		ConfigJSON: `{"devices":[]}`,
	}

	if err := db.SaveScene(scene); err != nil {
		t.Fatalf("SaveScene failed: %v", err)
	}

	scenes, err := db.GetScenes()
	if err != nil {
		t.Fatalf("GetScenes failed: %v", err)
	}
	if len(scenes) != 1 {
		t.Fatalf("Expected 1 scene, got %d", len(scenes))
	}

	// 3. Test DeleteGroup & DeleteScene
	if err := db.DeleteGroup("group-1"); err != nil {
		t.Fatalf("DeleteGroup failed: %v", err)
	}
	groups, _ = db.GetGroups()
	if len(groups) != 0 {
		t.Errorf("Expected 0 groups after delete, got %d", len(groups))
	}

	if err := db.DeleteScene("scene-1"); err != nil {
		t.Fatalf("DeleteScene failed: %v", err)
	}
	scenes, _ = db.GetScenes()
	if len(scenes) != 0 {
		t.Errorf("Expected 0 scenes after delete, got %d", len(scenes))
	}
}

func TestDatabase_DashboardItemsAndPanels(t *testing.T) {
	db := setupTestDB(t)

	// 1. Test SaveDashboardPanel & GetDashboardPanels
	panel := DashboardPanel{
		ID:       "panel-living-room",
		Title:    "Living Room Zone",
		Position: 0,
	}

	if err := db.SaveDashboardPanel(panel); err != nil {
		t.Fatalf("SaveDashboardPanel failed: %v", err)
	}

	panels, err := db.GetDashboardPanels()
	if err != nil {
		t.Fatalf("GetDashboardPanels failed: %v", err)
	}
	if len(panels) != 1 {
		t.Fatalf("Expected 1 panel, got %d", len(panels))
	}
	if panels[0].Title != "Living Room Zone" {
		t.Errorf("Expected panel title 'Living Room Zone', got '%s'", panels[0].Title)
	}

	// 2. Test UPSERT behavior for UpdateDashboardItemPanel on a new device
	if err := db.UpdateDashboardItemPanel("wled-dev-1", "panel-living-room"); err != nil {
		t.Fatalf("UpdateDashboardItemPanel failed: %v", err)
	}

	items, err := db.GetDashboardItems()
	if err != nil {
		t.Fatalf("GetDashboardItems failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("Expected 1 dashboard item after UPSERT, got %d", len(items))
	}
	if items[0].PanelID != "panel-living-room" {
		t.Errorf("Expected PanelID 'panel-living-room', got '%s'", items[0].PanelID)
	}

	// 3. Test UpdateDashboardItemSize UPSERT
	if err := db.UpdateDashboardItemSize("wled-dev-1", "wide"); err != nil {
		t.Fatalf("UpdateDashboardItemSize failed: %v", err)
	}

	items, _ = db.GetDashboardItems()
	if items[0].Size != "wide" {
		t.Errorf("Expected Size 'wide', got '%s'", items[0].Size)
	}

	// 4. Test PinDashboardItem UPSERT
	if err := db.PinDashboardItem("wled-dev-1", "device", false); err != nil {
		t.Fatalf("PinDashboardItem failed: %v", err)
	}

	items, _ = db.GetDashboardItems()
	if items[0].IsPinned != false {
		t.Errorf("Expected IsPinned false, got %v", items[0].IsPinned)
	}
	if items[0].PanelID != "panel-living-room" {
		t.Errorf("Expected PanelID 'panel-living-room' preserved, got '%s'", items[0].PanelID)
	}

	// 5. Test DeleteDashboardPanel unassigns item panel IDs
	if err := db.DeleteDashboardPanel("panel-living-room"); err != nil {
		t.Fatalf("DeleteDashboardPanel failed: %v", err)
	}

	panels, _ = db.GetDashboardPanels()
	if len(panels) != 0 {
		t.Errorf("Expected 0 panels after deletion, got %d", len(panels))
	}

	items, _ = db.GetDashboardItems()
	if items[0].PanelID != "" {
		t.Errorf("Expected PanelID to be reset to empty string after panel deletion, got '%s'", items[0].PanelID)
	}
}
