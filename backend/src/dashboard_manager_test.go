package main

import (
	"testing"
)

func TestDashboardManager_PanelRetention(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()

	dm := NewDashboardManager(db, hub)

	// 1. Assign a panel ID to device "wled-dev-10"
	item, err := dm.SetItemPanel("wled-dev-10", "panel-office")
	if err != nil {
		t.Fatalf("SetItemPanel failed: %v", err)
	}
	if item.PanelID != "panel-office" {
		t.Fatalf("Expected PanelID 'panel-office', got '%s'", item.PanelID)
	}

	// 2. Simulate mDNS discovery triggering AutoPinIfNew
	dm.AutoPinIfNew("wled-dev-10", "device")

	// 3. Verify that PanelID was NOT wiped back to empty string
	items := dm.GetItems()
	var found *DashboardItem
	for _, it := range items {
		if it.ItemID == "wled-dev-10" {
			copyIt := it
			found = &copyIt
			break
		}
	}

	if found == nil {
		t.Fatalf("Device 'wled-dev-10' not found in DashboardManager items")
	}
	if found.PanelID != "panel-office" {
		t.Errorf("AutoPinIfNew wiped PanelID! Expected 'panel-office', got '%s'", found.PanelID)
	}
}

func TestDashboardManager_PanelsCRUD(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()
	dm := NewDashboardManager(db, hub)

	// 1. Add panel
	panel, err := dm.AddPanel("panel-100", "Gaming Setup")
	if err != nil {
		t.Fatalf("AddPanel failed: %v", err)
	}
	if panel.Title != "Gaming Setup" {
		t.Errorf("Expected title 'Gaming Setup', got '%s'", panel.Title)
	}

	// 2. Get panels
	panels, err := dm.GetPanels()
	if err != nil {
		t.Fatalf("GetPanels failed: %v", err)
	}
	if len(panels) != 1 {
		t.Fatalf("Expected 1 panel, got %d", len(panels))
	}

	// 3. Delete panel
	if err := dm.DeletePanel("panel-100"); err != nil {
		t.Fatalf("DeletePanel failed: %v", err)
	}

	panels, _ = dm.GetPanels()
	if len(panels) != 0 {
		t.Errorf("Expected 0 panels after deletion, got %d", len(panels))
	}
}

func TestDashboardManager_Reorder(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()
	dm := NewDashboardManager(db, hub)

	dm.PinItem("dev-1", "device", true)
	dm.PinItem("dev-2", "device", true)

	if err := dm.ReorderItems([]string{"dev-2", "dev-1"}); err != nil {
		t.Fatalf("ReorderItems failed: %v", err)
	}

	items := dm.GetItems()
	for _, it := range items {
		if it.ItemID == "dev-2" && it.Position != 0 {
			t.Errorf("Expected dev-2 position 0, got %d", it.Position)
		}
		if it.ItemID == "dev-1" && it.Position != 1 {
			t.Errorf("Expected dev-1 position 1, got %d", it.Position)
		}
	}
}
