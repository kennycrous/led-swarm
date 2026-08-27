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
		RoomID:   "room-living",
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
	placements, err := db.GetCanvasPlacements("room-living")
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
	if err := db.DeleteCanvasPlacement("wled-canvas-1", "room-living"); err != nil {
		t.Fatalf("DeleteCanvasPlacement failed: %v", err)
	}

	placementsAfter, _ := db.GetCanvasPlacements("room-living")
	if len(placementsAfter) != 0 {
		t.Errorf("Expected 0 placements after delete, got %d", len(placementsAfter))
	}
}

func TestCanvasManager_RoomsAndPlacements(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()

	// Pre-save 2 devices
	_ = db.SaveDevice(Device{ID: "wled-c1", Name: "Strip 1", IPAddress: "192.168.1.201", LEDCount: 30, IsOnline: true})
	_ = db.SaveDevice(Device{ID: "wled-c2", Name: "Strip 2", IPAddress: "192.168.1.202", LEDCount: 60, IsOnline: true})

	cm := NewCanvasManager(db, hub)

	// 1. Create Canvas Room
	room, err := cm.CreateRoom("Living Room", "Main TV setup", 2000, 1200, []string{"wled-c1", "wled-c2"})
	if err != nil {
		t.Fatalf("CreateRoom failed: %v", err)
	}
	if room.Title != "Living Room" {
		t.Errorf("Expected 'Living Room', got '%s'", room.Title)
	}

	// 2. Batch Save Placements for Room
	p1 := CanvasPlacement{DeviceID: "wled-c1", RoomID: room.ID, PosX: 100, PosY: 100, Rotation: 0, Scale: 1, Geometry: "strip"}
	p2 := CanvasPlacement{DeviceID: "wled-c2", RoomID: room.ID, PosX: 300, PosY: 200, Rotation: 90, Scale: 1.5, Geometry: "matrix"}

	if err := cm.BatchSavePlacements(room.ID, []CanvasPlacement{p1, p2}); err != nil {
		t.Fatalf("BatchSavePlacements failed: %v", err)
	}

	placements := cm.GetPlacementsForRoom(room.ID)
	if len(placements) != 2 {
		t.Fatalf("Expected 2 placements in CanvasManager for room %s, got %d", room.ID, len(placements))
	}

	// 3. Get Rooms List
	rooms := cm.GetRooms()
	if len(rooms) == 0 {
		t.Errorf("Expected rooms list to contain at least 1 room")
	}

	// 4. Delete Room
	if err := cm.DeleteRoom(room.ID); err != nil {
		t.Fatalf("DeleteRoom failed: %v", err)
	}
}

func TestCanvasManager_PersistenceAcrossRestart(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()

	cm1 := NewCanvasManager(db, hub)
	room, err := cm1.CreateRoom("Bedroom", "Bed backlighting", 2000, 1200, []string{"wled-b1"})
	if err != nil {
		t.Fatalf("CreateRoom failed: %v", err)
	}

	p1 := CanvasPlacement{DeviceID: "wled-b1", RoomID: room.ID, PosX: 450, PosY: 350, Rotation: 180, Scale: 1, Geometry: "strip"}
	if err := cm1.BatchSavePlacements(room.ID, []CanvasPlacement{p1}); err != nil {
		t.Fatalf("BatchSavePlacements failed: %v", err)
	}

	// Now simulate server restart by creating a new CanvasManager connected to the same DB
	cm2 := NewCanvasManager(db, hub)

	rooms := cm2.GetRooms()
	if len(rooms) != 1 {
		t.Fatalf("Expected 1 room after restart, got %d", len(rooms))
	}
	if rooms[0].Title != "Bedroom" {
		t.Errorf("Expected room title 'Bedroom', got '%s'", rooms[0].Title)
	}

	placements := cm2.GetPlacementsForRoom(room.ID)
	if len(placements) != 1 {
		t.Fatalf("Expected 1 placement after restart for room %s, got %d", room.ID, len(placements))
	}
	if placements[0].PosX != 450 || placements[0].Rotation != 180 {
		t.Errorf("Unexpected placement data after restart: %+v", placements[0])
	}
}

func TestCanvasManager_UpdateRoomTitle(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()
	cm := NewCanvasManager(db, hub)

	room, err := cm.CreateRoom("Initial Name", "Description", 2000, 1200, nil)
	if err != nil {
		t.Fatalf("CreateRoom failed: %v", err)
	}

	updated, err := cm.UpdateRoomTitle(room.ID, "Renamed Room")
	if err != nil {
		t.Fatalf("UpdateRoomTitle failed: %v", err)
	}

	if updated.Title != "Renamed Room" {
		t.Errorf("Expected title 'Renamed Room', got '%s'", updated.Title)
	}

	// Verify persistence in DB after manager re-initialization
	cmRestarted := NewCanvasManager(db, hub)
	rooms := cmRestarted.GetRooms()
	if len(rooms) != 1 {
		t.Fatalf("Expected 1 room after restart, got %d", len(rooms))
	}
	if rooms[0].Title != "Renamed Room" {
		t.Errorf("Expected persisted title 'Renamed Room', got '%s'", rooms[0].Title)
	}
}
