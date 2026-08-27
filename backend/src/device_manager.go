package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type DeviceManager struct {
	db           *Database
	wledClient   *WLEDClient
	hub          *Hub
	dashboardMgr *DashboardManager
	devices      map[string]*Device
	effects      []string
	palettes     []string
	mu           sync.RWMutex
}

func NewDeviceManager(db *Database, wledClient *WLEDClient, hub *Hub) *DeviceManager {
	dm := &DeviceManager{
		db:         db,
		wledClient: wledClient,
		hub:        hub,
		devices:    make(map[string]*Device),
		effects:    make([]string, 0),
		palettes:   make([]string, 0),
	}

	// Load stored devices from DB on startup
	if err := dm.loadStoredDevices(); err != nil {
		log.Printf("[DeviceManager] Warning loading stored devices: %v", err)
	}

	return dm
}

func (dm *DeviceManager) SetDashboardManager(dashboardMgr *DashboardManager) {
	dm.dashboardMgr = dashboardMgr
}

func (dm *DeviceManager) loadStoredDevices() error {
	stored, err := dm.db.GetDevices()
	if err != nil {
		return err
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	for _, dev := range stored {
		d := dev
		dm.devices[d.ID] = &d
		// Connect WLED WebSocket in background
		go dm.connectDeviceWS(&d)
	}

	log.Printf("[DeviceManager] Loaded %d devices from database", len(stored))
	return nil
}

func (dm *DeviceManager) RegisterDevice(dev Device) {
	dm.mu.Lock()
	existing, ok := dm.devices[dev.ID]
	if ok && existing.Name != "" && dev.Name == "" {
		dev.Name = existing.Name
	}
	dm.devices[dev.ID] = &dev
	dm.mu.Unlock()

	// Save to DB
	if err := dm.db.SaveDevice(dev); err != nil {
		log.Printf("[DeviceManager] Failed to save device %s to DB: %v", dev.ID, err)
	}

	// Fetch effects and palettes if not cached yet
	go dm.fetchMetadataIfMissing(dev.IPAddress)

	// Connect WebSocket for live push updates
	go dm.connectDeviceWS(&dev)

	// Broadcast update via WS Hub
	dm.broadcastState(dev)
}

func (dm *DeviceManager) AddDeviceByIP(rawIP string) (*Device, error) {
	ip := strings.TrimSpace(rawIP)
	ip = strings.TrimPrefix(ip, "http://")
	ip = strings.TrimPrefix(ip, "https://")
	ip = strings.TrimSuffix(ip, "/")

	log.Printf("[DeviceManager] Attempting to add WLED device at IP: %s", ip)
	info, err := dm.wledClient.FetchDeviceInfo(ip)
	if err != nil {
		log.Printf("[DeviceManager] Failed to fetch device info for %s: %v", ip, err)
		return nil, fmt.Errorf("failed to reach WLED device at %s: %w", ip, err)
	}

	id := info.Mac
	if id == "" {
		id = strings.ReplaceAll(ip, ".", "-")
	}

	name := info.Name
	if name == "" {
		name = fmt.Sprintf("WLED-%s", id[len(id)-4:])
	}

	dev := Device{
		ID:         id,
		Name:       name,
		IPAddress:  ip,
		MACAddress: info.Mac,
		LEDCount:   info.Leds.Count,
		IsOnline:   true,
	}

	dm.RegisterDevice(dev)
	return &dev, nil
}

func (dm *DeviceManager) connectDeviceWS(dev *Device) {
	wsURL := fmt.Sprintf("ws://%s/ws", dev.IPAddress)
	dialer := websocket.Dialer{HandshakeTimeout: 3 * time.Second}

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		log.Printf("[DeviceManager] WS connect failed for %s (%s): %v, fallback to REST", dev.Name, dev.IPAddress, err)
		dm.updateOnlineState(dev.ID, false)
		return
	}
	defer conn.Close()

	dm.updateOnlineState(dev.ID, true)
	log.Printf("[DeviceManager] Connected WS to WLED %s (%s)", dev.Name, dev.IPAddress)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[DeviceManager] WS disconnected for %s: %v", dev.Name, err)
			dm.updateOnlineState(dev.ID, false)
			break
		}

		dm.handleWLEDWSMessage(dev.ID, message)
	}
}

func (dm *DeviceManager) handleWLEDWSMessage(deviceID string, message []byte) {
	var payload struct {
		State *WLEDState `json:"state"`
		Info  *WLEDInfo  `json:"info"`
	}

	if err := json.Unmarshal(message, &payload); err != nil {
		return
	}

	dm.mu.Lock()
	dev, ok := dm.devices[deviceID]
	if !ok {
		dm.mu.Unlock()
		return
	}

	if payload.State != nil {
		dev.State = *payload.State
	}
	if payload.Info != nil {
		dev.LEDCount = payload.Info.Leds.Count
		dev.IsOnline = true
	}
	devCopy := *dev
	dm.mu.Unlock()

	dm.broadcastState(devCopy)
}

func (dm *DeviceManager) updateOnlineState(deviceID string, isOnline bool) {
	dm.mu.Lock()
	dev, ok := dm.devices[deviceID]
	if ok {
		dev.IsOnline = isOnline
		devCopy := *dev
		dm.mu.Unlock()
		dm.broadcastState(devCopy)
		dm.db.SaveDevice(devCopy)
	} else {
		dm.mu.Unlock()
	}
}

func (dm *DeviceManager) fetchMetadataIfMissing(ip string) {
	dm.mu.RLock()
	hasMetadata := len(dm.effects) > 0 && len(dm.palettes) > 0
	dm.mu.RUnlock()

	if hasMetadata {
		return
	}

	url := fmt.Sprintf("http://%s/json/eff", ip)
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		var effs []string
		if err := json.NewDecoder(resp.Body).Decode(&effs); err == nil {
			dm.mu.Lock()
			dm.effects = effs
			dm.mu.Unlock()
			log.Printf("[DeviceManager] Cached %d WLED effect names", len(effs))
		}
		resp.Body.Close()
	}

	urlPal := fmt.Sprintf("http://%s/json/pal", ip)
	respPal, err := http.Get(urlPal)
	if err == nil && respPal.StatusCode == http.StatusOK {
		var pals []string
		if err := json.NewDecoder(respPal.Body).Decode(&pals); err == nil {
			dm.mu.Lock()
			dm.palettes = pals
			dm.mu.Unlock()
			log.Printf("[DeviceManager] Cached %d WLED palette names", len(pals))
		}
		respPal.Body.Close()
	}
}

func (dm *DeviceManager) GetAllDevices() []Device {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	result := make([]Device, 0, len(dm.devices))
	for _, dev := range dm.devices {
		result = append(result, *dev)
	}
	return result
}

func (dm *DeviceManager) UpdateDeviceName(id string, newName string) error {
	dm.mu.Lock()
	if dev, ok := dm.devices[id]; ok {
		dev.Name = newName
	}
	dm.mu.Unlock()

	return dm.db.UpdateDeviceName(id, newName)
}

func (dm *DeviceManager) GetEffects() []string {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.effects
}

func (dm *DeviceManager) GetPalettes() []string {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.palettes
}

func (dm *DeviceManager) broadcastState(dev Device) {
	if dm.hub != nil {
		dm.hub.BroadcastJSON(map[string]interface{}{
			"type":   "device_state",
			"device": dev,
		})
	}
}
