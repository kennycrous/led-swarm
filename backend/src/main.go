package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:dist
var assets embed.FS

func main() {
	serverMode := flag.Bool("server", false, "Run in headless server mode (for Docker / headless OS)")
	port := flag.Int("port", 8080, "Port to listen on in server mode")
	dbPath := flag.String("db", "led-swarm.db", "Path to SQLite database file")
	flag.Parse()

	log.Println("[Main] Initializing LED Swarm Orchestrator Engine...")

	// 1. Initialize SQLite Database (Pure-Go CGO-free)
	db, err := NewDatabase(*dbPath)
	if err != nil {
		log.Fatalf("[Main] Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 2. Initialize Core Subsystems
	wledClient := NewWLEDClient()
	hub := NewHub()
	go hub.Run()

	devMgr := NewDeviceManager(db, wledClient, hub)
	groupMgr := NewGroupManager(db, wledClient, devMgr, hub)
	dashboardMgr := NewDashboardManager(db, hub)
	canvasMgr := NewCanvasManager(db, hub)
	devMgr.SetDashboardManager(dashboardMgr)

	scanner := NewMDNSScanner(db, wledClient, devMgr)

	// Trigger initial mDNS scan
	go func() {
		if err := scanner.StartScan(context.Background()); err != nil {
			log.Printf("[Main] Initial mDNS scan error: %v", err)
		}
	}()

	// 3. Select Operating Mode (Server vs Desktop App)
	if *serverMode || os.Getenv("SERVER_MODE") == "1" {
		log.Println("[Main] Starting in Headless Server Mode...")
		srv := NewServer(db, wledClient, devMgr, groupMgr, dashboardMgr, canvasMgr, hub, scanner, assets)
		if err := srv.Start(*port); err != nil {
			log.Fatalf("[Main] Server failure: %v", err)
		}
	} else {
		log.Println("[Main] Launching Native Wails Desktop Application...")
		app := NewApp(db, wledClient, devMgr, groupMgr, dashboardMgr, canvasMgr, hub, scanner)

		err := wails.Run(&options.App{
			Title:  "LED Swarm Orchestrator",
			Width:  1280,
			Height: 800,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			BackgroundColour: &options.RGBA{R: 6, G: 9, B: 14, A: 255},
			OnStartup:        app.startup,
			Bind: []interface{}{
				app,
			},
		})

		if err != nil {
			log.Fatalf("[Main] Desktop application error: %v", err)
		}
	}
}
