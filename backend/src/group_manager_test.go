package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestGroupManager_GroupsAndScenes(t *testing.T) {
	ts1, mockState1 := createMockWLEDServerWithMAC(t, "AA:BB:CC:11:22:33")
	ts2, mockState2 := createMockWLEDServerWithMAC(t, "AA:BB:CC:44:55:66")

	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()
	devMgr := NewDeviceManager(db, wledClient, hub)

	host1 := strings.TrimPrefix(ts1.URL, "http://")
	host2 := strings.TrimPrefix(ts2.URL, "http://")

	dev1 := Device{ID: "wled-1", Name: "Strip 1", IPAddress: host1, MACAddress: "AA:BB:CC:11:22:33", LEDCount: 60, IsOnline: true}
	dev2 := Device{ID: "wled-2", Name: "Strip 2", IPAddress: host2, MACAddress: "AA:BB:CC:44:55:66", LEDCount: 60, IsOnline: true}

	devMgr.RegisterDevice(dev1)
	devMgr.RegisterDevice(dev2)

	time.Sleep(50 * time.Millisecond)

	gm := NewGroupManager(db, wledClient, devMgr, hub)

	// 1. Save Group
	group, err := gm.SaveGroup("", "Living Room All", "All TV and ceiling lights", []string{dev1.ID, dev2.ID})
	if err != nil {
		t.Fatalf("SaveGroup failed: %v", err)
	}
	if group.Name != "Living Room All" {
		t.Errorf("Expected group name 'Living Room All', got '%s'", group.Name)
	}

	groups := gm.GetGroups()
	if len(groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(groups))
	}

	// 2. SetGroupState (Concurrent Multi-Strip Dispatching)
	jsonState, _ := json.Marshal(WLEDState{On: true, Brightness: 240})
	if err := gm.SetGroupState(group.ID, json.RawMessage(jsonState)); err != nil {
		t.Fatalf("SetGroupState failed: %v", err)
	}
	if mockState1.Brightness != 240 || mockState2.Brightness != 240 {
		t.Errorf("SetGroupState failed to dispatch to all group devices concurrently")
	}

	// 3. Capture Scene Snapshot
	scene, err := gm.CaptureScene("Movie Night", "Film")
	if err != nil {
		t.Fatalf("CaptureScene failed: %v", err)
	}
	if scene.Name != "Movie Night" {
		t.Errorf("Expected scene name 'Movie Night', got '%s'", scene.Name)
	}

	scenes := gm.GetScenes()
	if len(scenes) != 1 {
		t.Fatalf("Expected 1 scene, got %d", len(scenes))
	}

	// Mutate mock states
	mockState1.Brightness = 10
	mockState2.Brightness = 10

	// 4. ApplyScene (Restore Snapshot)
	if err := gm.ApplyScene(scene.ID); err != nil {
		t.Fatalf("ApplyScene failed: %v", err)
	}
	if mockState1.Brightness != 240 || mockState2.Brightness != 240 {
		t.Errorf("ApplyScene failed to restore state snapshot")
	}

	// 5. Delete Group & Scene
	if err := gm.DeleteGroup(group.ID); err != nil {
		t.Fatalf("DeleteGroup failed: %v", err)
	}
	if len(gm.GetGroups()) != 0 {
		t.Errorf("Expected 0 groups after delete")
	}

	if err := gm.DeleteScene(scene.ID); err != nil {
		t.Fatalf("DeleteScene failed: %v", err)
	}
	if len(gm.GetScenes()) != 0 {
		t.Errorf("Expected 0 scenes after delete")
	}
}

func TestGroupManager_ScopedSceneRestore(t *testing.T) {
	ts1, mockState1 := createMockWLEDServerWithMAC(t, "AA:BB:CC:77:88:99")
	ts2, mockState2 := createMockWLEDServerWithMAC(t, "AA:BB:CC:00:11:22")

	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()
	devMgr := NewDeviceManager(db, wledClient, hub)

	host1 := strings.TrimPrefix(ts1.URL, "http://")
	host2 := strings.TrimPrefix(ts2.URL, "http://")

	dev1 := Device{ID: "wled-scoped-1", Name: "Office Strip", IPAddress: host1, MACAddress: "AA:BB:CC:77:88:99", LEDCount: 60, IsOnline: true}
	dev2 := Device{ID: "wled-scoped-2", Name: "Bedroom Strip", IPAddress: host2, MACAddress: "AA:BB:CC:00:11:22", LEDCount: 60, IsOnline: true}

	devMgr.RegisterDevice(dev1)
	devMgr.RegisterDevice(dev2)

	time.Sleep(50 * time.Millisecond)

	gm := NewGroupManager(db, wledClient, devMgr, hub)

	// Create Group for dev1 only
	group, _ := gm.SaveGroup("", "Office Only", "Office Desk", []string{dev1.ID})

	// Set initial states: Office = 200, Bedroom = 50
	mockState1.Brightness = 200
	mockState2.Brightness = 50

	// Capture Scoped Scene for group
	scene, err := gm.CaptureScopedScene("Office Bright", "Sun", "group", group.ID)
	if err != nil {
		t.Fatalf("CaptureScopedScene failed: %v", err)
	}
	if scene.ScopeType != "group" || scene.TargetID != group.ID {
		t.Errorf("Unexpected scope configuration: scopeType=%s, targetId=%s", scene.ScopeType, scene.TargetID)
	}

	// Mutate states
	mockState1.Brightness = 10
	mockState2.Brightness = 99

	// Apply Scoped Scene -> Should ONLY restore dev1 (Office) to 200, leaving dev2 (Bedroom) untouched at 99!
	if err := gm.ApplyScene(scene.ID); err != nil {
		t.Fatalf("ApplyScene failed for scoped scene: %v", err)
	}

	if mockState1.Brightness != 200 {
		t.Errorf("Expected Office strip to be restored to 200, got %d", mockState1.Brightness)
	}
	if mockState2.Brightness != 99 {
		t.Errorf("Expected Bedroom strip to remain untouched at 99, got %d", mockState2.Brightness)
	}
}
