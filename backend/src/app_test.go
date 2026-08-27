package main

import (
	"context"
	"strings"
	"testing"
)

func TestApp_WailsBindings(t *testing.T) {
	ts, mockState := createMockWLEDServer(t)
	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()

	devMgr := NewDeviceManager(db, wledClient, hub)
	groupMgr := NewGroupManager(db, wledClient, devMgr, hub)
	dashboardMgr := NewDashboardManager(db, hub)
	canvasMgr := NewCanvasManager(db, hub)
	ddpStreamer := NewDDPStreamer()
	devMgr.SetDashboardManager(dashboardMgr)

	app := NewApp(db, wledClient, devMgr, groupMgr, dashboardMgr, canvasMgr, ddpStreamer, hub, nil)
	app.startup(context.Background())

	host := strings.TrimPrefix(ts.URL, "http://")

	// 1. AddManualDevice
	dev, err := app.AddManualDevice(host)
	if err != nil {
		t.Fatalf("App.AddManualDevice failed: %v", err)
	}
	if dev.Name != "Mock WLED Strip" {
		t.Errorf("Expected 'Mock WLED Strip', got '%s'", dev.Name)
	}

	// 2. GetDevices
	devices, err := app.GetDevices()
	if err != nil || len(devices) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(devices))
	}

	// 3. SetDevicePower / SetDeviceBrightness / SetDeviceColor / SetDeviceEffect / SetDevicePalette
	if err := app.SetDevicePower(host, false); err != nil {
		t.Fatalf("App.SetDevicePower failed: %v", err)
	}
	if mockState.On != false {
		t.Errorf("Expected mockState.On=false")
	}

	if err := app.SetDeviceBrightness(host, 150); err != nil {
		t.Fatalf("App.SetDeviceBrightness failed: %v", err)
	}
	if err := app.SetDeviceColor(host, 255, 0, 0); err != nil {
		t.Fatalf("App.SetDeviceColor failed: %v", err)
	}
	if err := app.SetDeviceEffect(host, 1); err != nil {
		t.Fatalf("App.SetDeviceEffect failed: %v", err)
	}
	if err := app.SetDevicePalette(host, 2); err != nil {
		t.Fatalf("App.SetDevicePalette failed: %v", err)
	}

	// 4. UpdateDeviceNickname
	if err := app.UpdateDeviceNickname(dev.ID, "Renamed Strip"); err != nil {
		t.Fatalf("App.UpdateDeviceNickname failed: %v", err)
	}

	// 5. GetEffects & GetPalettes
	effects, _ := app.GetEffects()
	palettes, _ := app.GetPalettes()
	if len(effects) == 0 || len(palettes) == 0 {
		t.Errorf("Expected non-empty effects/palettes array")
	}

	// 6. SaveGroup & GetGroups
	group, err := app.SaveGroup("", "Office Zone", "Desk lights", []string{dev.ID})
	if err != nil {
		t.Fatalf("App.SaveGroup failed: %v", err)
	}
	groups, _ := app.GetGroups()
	if len(groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(groups))
	}

	// 7. CaptureScene & GetScenes & ApplyScene
	scene, err := app.CaptureScene("Focus Mode", "Zap")
	if err != nil {
		t.Fatalf("App.CaptureScene failed: %v", err)
	}
	scenes, _ := app.GetScenes()
	if len(scenes) != 1 {
		t.Fatalf("Expected 1 scene, got %d", len(scenes))
	}
	if err := app.ApplyScene(scene.ID); err != nil {
		t.Fatalf("App.ApplyScene failed: %v", err)
	}

	// 8. AddDashboardPanel & GetDashboardPanels & PinDashboardItem & SetDashboardItemSize & SetDashboardItemPanel
	panel, err := app.AddDashboardPanel("Desk Panel")
	if err != nil {
		t.Fatalf("App.AddDashboardPanel failed: %v", err)
	}
	panels, err := app.GetDashboardPanels()
	if err != nil || len(panels) != 1 {
		t.Fatalf("Expected 1 panel from App binding, got %d", len(panels))
	}

	if _, err := app.PinDashboardItem(dev.ID, "device", true); err != nil {
		t.Fatalf("App.PinDashboardItem failed: %v", err)
	}
	if _, err := app.SetDashboardItemSize(dev.ID, "wide"); err != nil {
		t.Fatalf("App.SetDashboardItemSize failed: %v", err)
	}
	if _, err := app.SetDashboardItemPanel(dev.ID, panel.ID); err != nil {
		t.Fatalf("App.SetDashboardItemPanel failed: %v", err)
	}
	items, _ := app.GetDashboardItems()
	if len(items) != 1 {
		t.Fatalf("Expected 1 dashboard item, got %d", len(items))
	}

	// 9. Delete operations
	if err := app.DeleteGroup(group.ID); err != nil {
		t.Fatalf("App.DeleteGroup failed: %v", err)
	}
	if err := app.DeleteScene(scene.ID); err != nil {
		t.Fatalf("App.DeleteScene failed: %v", err)
	}
	if err := app.DeleteDashboardPanel(panel.ID); err != nil {
		t.Fatalf("App.DeleteDashboardPanel failed: %v", err)
	}
}
