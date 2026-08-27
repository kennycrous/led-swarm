package main

import (
	"context"
	"log"
)

// App struct for Wails desktop application
type App struct {
	ctx        context.Context
	db         *Database
	wledClient *WLEDClient
	hub        *Hub
	scanner    *MDNSScanner
}

func NewApp(db *Database, wledClient *WLEDClient, hub *Hub, scanner *MDNSScanner) *App {
	return &App{
		db:         db,
		wledClient: wledClient,
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
	return a.db.GetDevices()
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
