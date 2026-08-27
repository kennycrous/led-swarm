package main

import (
	"context"
	"encoding/json"
	"log"
)

// App struct for Wails desktop application
type App struct {
	ctx        context.Context
	db         *Database
	wledClient *WLEDClient
	devMgr     *DeviceManager
	groupMgr   *GroupManager
	hub        *Hub
	scanner    *MDNSScanner
}

func NewApp(db *Database, wledClient *WLEDClient, devMgr *DeviceManager, groupMgr *GroupManager, hub *Hub, scanner *MDNSScanner) *App {
	return &App{
		db:         db,
		wledClient: wledClient,
		devMgr:     devMgr,
		groupMgr:   groupMgr,
		hub:        hub,
		scanner:    scanner,
	}
}

// startup is called when Wails desktop app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("[Wails App] Desktop application initialized")
}

// Wails Exposed API Methods (Auto-generated TypeScript bindings)

func (a *App) GetDevices() ([]Device, error) {
	return a.devMgr.GetAllDevices(), nil
}

func (a *App) TriggerScan() string {
	go func() {
		if err := a.scanner.StartScan(a.ctx); err != nil {
			log.Printf("[Wails App] Scan error: %v", err)
		}
	}()
	return "scan_started"
}

func (a *App) SetDevicePower(ip string, on bool) error {
	state := WLEDState{On: on}
	return a.wledClient.SetState(ip, state)
}

func (a *App) SetDeviceBrightness(ip string, brightness int) error {
	state := WLEDState{Brightness: brightness}
	return a.wledClient.SetState(ip, state)
}

func (a *App) SetDeviceColor(ip string, r int, g int, b int) error {
	segment := WLEDSegment{
		ID:     0,
		Colors: [][]int{{r, g, b}},
	}
	state := WLEDState{Segments: []WLEDSegment{segment}}
	return a.wledClient.SetState(ip, state)
}

func (a *App) SetDeviceEffect(ip string, fx int) error {
	segment := WLEDSegment{
		ID: 0,
		FX: fx,
	}
	state := WLEDState{Segments: []WLEDSegment{segment}}
	return a.wledClient.SetState(ip, state)
}

func (a *App) SetDevicePalette(ip string, pal int) error {
	segment := WLEDSegment{
		ID:  0,
		Pal: pal,
	}
	state := WLEDState{Segments: []WLEDSegment{segment}}
	return a.wledClient.SetState(ip, state)
}

func (a *App) AddManualDevice(ip string) (*Device, error) {
	return a.devMgr.AddDeviceByIP(ip)
}

func (a *App) UpdateDeviceNickname(id string, name string) error {
	return a.devMgr.UpdateDeviceName(id, name)
}

func (a *App) GetEffects() ([]string, error) {
	return a.devMgr.GetEffects(), nil
}

func (a *App) GetPalettes() ([]string, error) {
	return a.devMgr.GetPalettes(), nil
}

// Group Wails Bindings

func (a *App) GetGroups() ([]Group, error) {
	return a.groupMgr.GetGroups(), nil
}

func (a *App) SaveGroup(id string, name string, description string, deviceIDs []string) (*Group, error) {
	return a.groupMgr.SaveGroup(id, name, description, deviceIDs)
}

func (a *App) DeleteGroup(id string) error {
	return a.groupMgr.DeleteGroup(id)
}

func (a *App) SetGroupState(groupID string, rawStateJSON string) error {
	return a.groupMgr.SetGroupState(groupID, json.RawMessage(rawStateJSON))
}

// Scene Wails Bindings

func (a *App) GetScenes() ([]Scene, error) {
	return a.groupMgr.GetScenes(), nil
}

func (a *App) CaptureScene(name string, icon string) (*Scene, error) {
	return a.groupMgr.CaptureScene(name, icon)
}

func (a *App) ApplyScene(id string) error {
	return a.groupMgr.ApplyScene(id)
}

func (a *App) DeleteScene(id string) error {
	return a.groupMgr.DeleteScene(id)
}
