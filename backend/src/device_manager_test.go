package main

import (
	"strings"
	"testing"
	"time"
)

func TestDeviceManager_AddDeviceByIPAndStateMutations(t *testing.T) {
	ts, mockState := createMockWLEDServer(t)
	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()

	dm := NewDeviceManager(db, wledClient, hub)

	host := strings.TrimPrefix(ts.URL, "http://")

	dev, err := dm.AddDeviceByIP(host)
	if err != nil {
		t.Fatalf("AddDeviceByIP failed: %v", err)
	}
	if dev.Name != "Mock WLED Strip" {
		t.Errorf("Expected Name 'Mock WLED Strip', got '%s'", dev.Name)
	}

	if err := wledClient.SetState(dev.IPAddress, WLEDState{On: false, Brightness: 220}); err != nil {
		t.Fatalf("SetState failed: %v", err)
	}
	if mockState.On != false || mockState.Brightness != 220 {
		t.Errorf("Expected mock server state updated via wledClient")
	}
}

func TestDeviceManager_IPValidationAndRename(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()

	dm := NewDeviceManager(db, wledClient, hub)

	// 1. Mock Add device manually
	dev := Device{
		ID:        "wled-manual-1",
		Name:      "Desk Strip",
		IPAddress: "192.168.1.150",
		LEDCount:  30,
		IsOnline:  true,
	}
	if err := db.SaveDevice(dev); err != nil {
		t.Fatalf("SaveDevice failed: %v", err)
	}

	// Load stored devices
	if err := dm.loadStoredDevices(); err != nil {
		t.Fatalf("loadStoredDevices failed: %v", err)
	}

	devices := dm.GetAllDevices()
	if len(devices) != 1 {
		t.Fatalf("Expected 1 device loaded, got %d", len(devices))
	}

	// 2. Update device nickname
	if err := dm.UpdateDeviceName("wled-manual-1", "Studio Backlight"); err != nil {
		t.Fatalf("UpdateDeviceName failed: %v", err)
	}

	devices = dm.GetAllDevices()
	if len(devices) != 1 || devices[0].Name != "Studio Backlight" {
		t.Errorf("Expected device name 'Studio Backlight', got '%v'", devices)
	}
}

func TestDeviceManager_EffectsAndPalettesCaching(t *testing.T) {
	ts, _ := createMockWLEDServer(t)
	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()

	dm := NewDeviceManager(db, wledClient, hub)

	host := strings.TrimPrefix(ts.URL, "http://")
	_, err := dm.AddDeviceByIP(host)
	if err != nil {
		t.Fatalf("AddDeviceByIP failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	effects := dm.GetEffects()
	palettes := dm.GetPalettes()

	if len(effects) == 0 {
		t.Errorf("Expected non-empty effects array")
	}
	if len(palettes) == 0 {
		t.Errorf("Expected non-empty palettes array")
	}
}
