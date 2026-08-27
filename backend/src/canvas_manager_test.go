package main

import (
	"testing"
)

func TestDatabase_CanvasPlacements(t *testing.T) {
	db := setupTestDB(t)

	// Save test device first (for foreign key integrity)
	dev := Device{
		ID:        "wled-canvas-1",
		Name:      "Room Strip 1",
		IPAddress: "192.168.1.101",
		LEDCount:  60,
		IsOnline:  true,
	}
	if err := db.SaveDevice(dev); err != nil {
		t.Fatalf("SaveDevice failed: %v", err)
	}

	placement := CanvasPlacement{
		DeviceID: "wled-canvas-1",
		PosX:     150.5,
		PosY:     250.0,
		Rotation: 45.0,
		Scale:    1.2,
		Geometry: "strip",
	}

	// 1. SaveCanvasPlacement
	if err := db.SaveCanvasPlacement(placement); err != nil {
		t.Fatalf("SaveCanvasPlacement failed: %v", err)
	}

	// 2. GetCanvasPlacements
	placements, err := db.GetCanvasPlacements()
	if err != nil {
		t.Fatalf("GetCanvasPlacements failed: %v", err)
	}
	if len(placements) != 1 {
		t.Fatalf("Expected 1 placement, got %d", len(placements))
	}
	if placements[0].PosX != 150.5 || placements[0].Rotation != 45.0 {
		t.Errorf("Unexpected placement data: %+v", placements[0])
	}

	// 3. DeleteCanvasPlacement
	if err := db.DeleteCanvasPlacement("wled-canvas-1"); err != nil {
		t.Fatalf("DeleteCanvasPlacement failed: %v", err)
	}

	placementsAfter, _ := db.GetCanvasPlacements()
	if len(placementsAfter) != 0 {
		t.Errorf("Expected 0 placements after delete, got %d", len(placementsAfter))
	}
}

func TestCanvasManager_PlacementsAndBatchSave(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()

	// Pre-save 2 devices
	_ = db.SaveDevice(Device{ID: "wled-c1", Name: "Strip 1", IPAddress: "192.168.1.201", LEDCount: 30, IsOnline: true})
	_ = db.SaveDevice(Device{ID: "wled-c2", Name: "Strip 2", IPAddress: "192.168.1.202", LEDCount: 60, IsOnline: true})

	cm := NewCanvasManager(db, hub)

	p1 := CanvasPlacement{DeviceID: "wled-c1", PosX: 100, PosY: 100, Rotation: 0, Scale: 1, Geometry: "strip"}
	p2 := CanvasPlacement{DeviceID: "wled-c2", PosX: 300, PosY: 200, Rotation: 90, Scale: 1.5, Geometry: "matrix"}

	if err := cm.BatchSavePlacements([]CanvasPlacement{p1, p2}); err != nil {
		t.Fatalf("BatchSavePlacements failed: %v", err)
	}

	placements := cm.GetPlacements()
	if len(placements) != 2 {
		t.Fatalf("Expected 2 placements in CanvasManager, got %d", len(placements))
	}

	pMap := make(map[string]CanvasPlacement)
	for _, p := range placements {
		pMap[p.DeviceID] = p
	}

	if pMap["wled-c2"].Rotation != 90 || pMap["wled-c2"].Geometry != "matrix" {
		t.Errorf("Unexpected CanvasPlacement for wled-c2: %+v", pMap["wled-c2"])
	}
}
