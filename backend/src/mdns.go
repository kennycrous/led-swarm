package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

type MDNSScanner struct {
	db         *Database
	wledClient *WLEDClient
	hub        *Hub
	isScanning bool
	mu         sync.Mutex
}

func NewMDNSScanner(db *Database, wledClient *WLEDClient, hub *Hub) *MDNSScanner {
	return &MDNSScanner{
		db:         db,
		wledClient: wledClient,
		hub:        hub,
	}
}

func (s *MDNSScanner) StartScan(ctx context.Context) error {
	s.mu.Lock()
	if s.isScanning {
		s.mu.Unlock()
		return nil
	}
	s.isScanning = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.isScanning = false
		s.mu.Unlock()
	}()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Printf("[mDNS] Failed to initialize resolver: %v", err)
		return err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			s.handleServiceEntry(entry)
		}
	}(entries)

	scanCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = resolver.Browse(scanCtx, "_wled._tcp", "local.", entries)
	if err != nil {
		log.Printf("[mDNS] Failed to browse service _wled._tcp: %v", err)
		return err
	}

	<-scanCtx.Done()
	return nil
}

func (s *MDNSScanner) handleServiceEntry(entry *zeroconf.ServiceEntry) {
	if len(entry.AddrIPv4) == 0 {
		return
	}

	ip := entry.AddrIPv4[0].String()
	name := strings.TrimSuffix(entry.Instance, ".local")

	info, err := s.wledClient.FetchDeviceInfo(ip)
	if err != nil {
		log.Printf("[mDNS] Discovered service at %s, but failed to fetch info: %v", ip, err)
		return
	}

	id := info.Mac
	if id == "" {
		id = strings.ReplaceAll(ip, ".", "-")
	}

	dev := Device{
		ID:         id,
		Name:       name,
		IPAddress:  ip,
		MACAddress: info.Mac,
		LEDCount:   info.Leds.Count,
		IsOnline:   true,
	}

	if err := s.db.SaveDevice(dev); err != nil {
		log.Printf("[mDNS] Failed to save device %s to database: %v", id, err)
		return
	}

	log.Printf("[mDNS] Successfully registered WLED device: %s (%s) [%d LEDs]", dev.Name, dev.IPAddress, dev.LEDCount)

	// Broadcast update to WebSocket hub
	if s.hub != nil {
		s.hub.BroadcastJSON(map[string]interface{}{
			"type":   "device_discovered",
			"device": dev,
		})
	}
}
