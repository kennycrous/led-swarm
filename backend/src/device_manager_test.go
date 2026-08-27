package main

import (
	"testing"
)

func TestDeviceManager_IPValidationAndRename(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()

	dm := NewDeviceManager(db, wledClient, hub)

	// 1. Invalid IP validation
	_, err := dm.AddDeviceByIP("invalid-ip")
	if err == nil {
		t.Errorf("Expected error for invalid IP string, got nil")
	}

	// 2. Mock Add device manually
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

	// 3. Update device nickname
	if err := dm.UpdateDeviceName("wled-manual-1", "Studio Backlight"); err != nil {
		t.Fatalf("UpdateDeviceName failed: %v", err)
	}

	devices = dm.GetAllDevices()
	if len(devices) != 1 || devices[0].Name != "Studio Backlight" {
		t.Errorf("Expected device name 'Studio Backlight', got '%v'", devices)
	}
}

func TestDeviceManager_EffectsAndPalettesCaching(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()

	dm := NewDeviceManager(db, wledClient, hub)

	effects := dm.GetEffects()
	palettes := dm.GetPalettes()

	if len(effects) == 0 {
		t.Errorf("Expected non-empty effects array fallback")
	}
	if len(palettes) == 0 {
		t.Errorf("Expected non-empty palettes array fallback")
	}
}
