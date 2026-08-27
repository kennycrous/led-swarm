package main

import (
	"testing"
)

func TestGroupManager_GroupsAndScenes(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()
	devMgr := NewDeviceManager(db, wledClient, hub)

	gm := NewGroupManager(db, wledClient, devMgr, hub)

	// 1. Save Group
	group, err := gm.SaveGroup("", "Living Room All", "All TV and ceiling lights", []string{"wled-1", "wled-2"})
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

	// 2. Capture Scene Snapshot
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

	// 3. Delete Group & Scene
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
