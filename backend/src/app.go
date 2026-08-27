package main

import (
	"context"
	"encoding/json"
	"log"
)

// App struct for Wails desktop application
type App struct {
	ctx          context.Context
	db           *Database
	wledClient   *WLEDClient
	devMgr       *DeviceManager
	groupMgr     *GroupManager
	dashboardMgr *DashboardManager
	canvasMgr    *CanvasManager
	hub          *Hub
	scanner      *MDNSScanner
}

func NewApp(db *Database, wledClient *WLEDClient, devMgr *DeviceManager, groupMgr *GroupManager, dashboardMgr *DashboardManager, canvasMgr *CanvasManager, hub *Hub, scanner *MDNSScanner) *App {
	return &App{
		db:           db,
		wledClient:   wledClient,
		devMgr:       devMgr,
		groupMgr:     groupMgr,
		dashboardMgr: dashboardMgr,
		canvasMgr:    canvasMgr,
		hub:          hub,
		scanner:      scanner,
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

func (a *App) CaptureScopedScene(name string, icon string, scopeType string, targetID string) (*Scene, error) {
	return a.groupMgr.CaptureScopedScene(name, icon, scopeType, targetID)
}

func (a *App) ApplyScene(id string) error {
	return a.groupMgr.ApplyScene(id)
}

func (a *App) DeleteScene(id string) error {
	return a.groupMgr.DeleteScene(id)
}

// Dashboard Wails Bindings

func (a *App) GetDashboardItems() ([]DashboardItem, error) {
	return a.dashboardMgr.GetItems(), nil
}

func (a *App) PinDashboardItem(itemID string, itemType string, isPinned bool) (*DashboardItem, error) {
	return a.dashboardMgr.PinItem(itemID, itemType, isPinned)
}

func (a *App) SetDashboardItemSize(itemID string, size string) (*DashboardItem, error) {
	return a.dashboardMgr.SetItemSize(itemID, size)
}

func (a *App) SetDashboardItemPanel(itemID string, panelID string) (*DashboardItem, error) {
	return a.dashboardMgr.SetItemPanel(itemID, panelID)
}

func (a *App) GetDashboardPanels() ([]DashboardPanel, error) {
	return a.dashboardMgr.GetPanels()
}

func (a *App) UpdateGroupName(id string, name string) error {
	_, err := a.groupMgr.RenameGroup(id, name)
	return err
}

func (a *App) AddDashboardPanel(title string) (*DashboardPanel, error) {
	return a.dashboardMgr.AddPanel("", title)
}

func (a *App) RenameDashboardPanel(id string, title string) (*DashboardPanel, error) {
	return a.dashboardMgr.RenamePanel(id, title)
}

func (a *App) DeleteDashboardPanel(id string) error {
	return a.dashboardMgr.DeletePanel(id)
}

func (a *App) GetCanvasRooms() ([]CanvasRoom, error) {
	return a.canvasMgr.GetRooms(), nil
}

func (a *App) CreateCanvasRoom(title string, description string, width int, height int, deviceIDs []string) (CanvasRoom, error) {
	return a.canvasMgr.CreateRoom(title, description, width, height, deviceIDs)
}

func (a *App) UpdateRoomTitle(id string, title string) (*CanvasRoom, error) {
	return a.canvasMgr.UpdateRoomTitle(id, title)
}

func (a *App) DeleteCanvasRoom(id string) error {
	return a.canvasMgr.DeleteRoom(id)
}

func (a *App) GetCanvasPlacements(roomID string) ([]CanvasPlacement, error) {
	return a.canvasMgr.GetPlacementsForRoom(roomID), nil
}

func (a *App) SaveCanvasPlacement(placement CanvasPlacement) error {
	return a.canvasMgr.SavePlacement(placement)
}

func (a *App) BatchSaveCanvasPlacements(roomID string, placements []CanvasPlacement) error {
	return a.canvasMgr.BatchSavePlacements(roomID, placements)
}
