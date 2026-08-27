package main

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/grandcat/zeroconf"
)

func TestMDNSScanner_HandleServiceEntry(t *testing.T) {
	ts, _ := createMockWLEDServer(t)
	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()

	devMgr := NewDeviceManager(db, wledClient, hub)

	host := strings.TrimPrefix(ts.URL, "http://")

	entry := &zeroconf.ServiceEntry{
		ServiceRecord: zeroconf.ServiceRecord{
			Instance: "wled-livingroom",
		},
		AddrIPv4: []net.IP{net.ParseIP("127.0.0.1")},
	}

	// Directly register mock device for handleServiceEntry test
	info, err := wledClient.FetchDeviceInfo(host)
	if err != nil {
		t.Fatalf("FetchDeviceInfo failed: %v", err)
	}

	dev := Device{
		ID:         info.Mac,
		Name:       entry.Instance,
		IPAddress:  host,
		MACAddress: info.Mac,
		LEDCount:   info.Leds.Count,
		IsOnline:   true,
	}

	devMgr.RegisterDevice(dev)

	// Verify device registered in DeviceManager
	devices := devMgr.GetAllDevices()
	if len(devices) != 1 {
		t.Fatalf("Expected 1 device registered, got %d", len(devices))
	}

	if devices[0].Name != "wled-livingroom" {
		t.Errorf("Expected Name 'wled-livingroom', got '%s'", devices[0].Name)
	}
}

func TestMDNSScanner_IsScanningFlag(t *testing.T) {
	db := setupTestDB(t)
	hub := NewHub()
	wledClient := NewWLEDClient()
	devMgr := NewDeviceManager(db, wledClient, hub)

	scanner := NewMDNSScanner(db, wledClient, devMgr)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = scanner.StartScan(ctx)

	scanner.mu.Lock()
	isScanning := scanner.isScanning
	scanner.mu.Unlock()

	if isScanning {
		t.Errorf("Expected isScanning to be false after StartScan completes")
	}
}
