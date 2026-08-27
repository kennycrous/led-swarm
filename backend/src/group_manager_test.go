package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGroupManager_GroupsAndScenes(t *testing.T) {
	ts1, mockState1 := createMockWLEDServer(t)
	ts2, mockState2 := createMockWLEDServer(t)

	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()
	devMgr := NewDeviceManager(db, wledClient, hub)

	host1 := strings.TrimPrefix(ts1.URL, "http://")
	host2 := strings.TrimPrefix(ts2.URL, "http://")

	dev1, _ := devMgr.AddDeviceByIP(host1)
	dev2, _ := devMgr.AddDeviceByIP(host2)

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
