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
	scanner := NewMDNSScanner(db, wledClient, devMgr)

	host := strings.TrimPrefix(ts.URL, "http://")
	ipStr, _, _ := net.SplitHostPort(host)
	if ipStr == "" {
		ipStr = "127.0.0.1"
	}

	entry := &zeroconf.ServiceEntry{
		ServiceRecord: zeroconf.ServiceRecord{
			Instance: "wled-livingroom",
		},
		AddrIPv4: []net.IP{net.ParseIP(ipStr)},
	}

	// Test handleServiceEntry
	scanner.handleServiceEntry(entry)

	// Verify device registered in DeviceManager
	devices := devMgr.GetAllDevices()
	if len(devices) != 1 {
		t.Fatalf("Expected 1 device registered via handleServiceEntry, got %d", len(devices))
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

	// Run scan (will timeout quickly since local network zeroconf browse returns in 100ms)
	_ = scanner.StartScan(ctx)

	scanner.mu.Lock()
	isScanning := scanner.isScanning
	scanner.mu.Unlock()

	if isScanning {
		t.Errorf("Expected isScanning to be false after StartScan completes")
	}
}
