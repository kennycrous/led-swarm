package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

// Server handles headless HTTP web API and static file serving
type Server struct {
	db         *Database
	wledClient *WLEDClient
	devMgr     *DeviceManager
	hub        *Hub
	scanner    *MDNSScanner
	distFS     embed.FS
}

func NewServer(db *Database, wledClient *WLEDClient, devMgr *DeviceManager, hub *Hub, scanner *MDNSScanner, distFS embed.FS) *Server {
	return &Server{
		db:         db,
		wledClient: wledClient,
		devMgr:     devMgr,
		hub:        hub,
		scanner:    scanner,
		distFS:     distFS,
	}
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()

	// REST API Endpoints
	mux.HandleFunc("/api/v1/devices", s.handleDevices)
	mux.HandleFunc("/api/v1/devices/add", s.handleAddDevice)
	mux.HandleFunc("/api/v1/devices/name", s.handleUpdateName)
	mux.HandleFunc("/api/v1/devices/state", s.handleSetState)
	mux.HandleFunc("/api/v1/effects", s.handleEffects)
	mux.HandleFunc("/api/v1/palettes", s.handlePalettes)
	mux.HandleFunc("/api/v1/scan", s.handleScan)
	mux.HandleFunc("/api/v1/ws", s.hub.ServeWS)

	// Embedded Frontend Static Asset File Server
	subFS, err := fs.Sub(s.distFS, "dist")
	if err != nil {
		log.Printf("[Server] Warning: 'dist' folder not found in embedded FS, serving API only")
	} else {
		fileServer := http.FileServer(http.FS(subFS))
		mux.Handle("/", fileServer)
	}

	addr := fmt.Sprintf(":%d", port)
	log.Printf("[Server] LED Swarm Orchestrator running at http://localhost%s", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	devices := s.devMgr.GetAllDevices()
	json.NewEncoder(w).Encode(devices)
}

func (s *Server) handleAddDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		IP string `json:"ip"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IP == "" {
		http.Error(w, "Invalid IP address provided", http.StatusBadRequest)
		return
	}

	dev, err := s.devMgr.AddDeviceByIP(req.IP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(dev)
}

func (s *Server) handleUpdateName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" || req.Name == "" {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	if err := s.devMgr.UpdateDeviceName(req.ID, req.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "name_updated"})
}

func (s *Server) handleSetState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		IP    string          `json:"ip"`
		State json.RawMessage `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IP == "" {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	if err := s.wledClient.SetRawState(req.IP, req.State); err != nil {
		log.Printf("[Server] Error updating WLED state for %s: %v", req.IP, err)
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "state_updated"})
}

func (s *Server) handleEffects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	effects := s.devMgr.GetEffects()
	json.NewEncoder(w).Encode(effects)
}

func (s *Server) handlePalettes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	palettes := s.devMgr.GetPalettes()
	json.NewEncoder(w).Encode(palettes)
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	go func() {
		if err := s.scanner.StartScan(r.Context()); err != nil {
			log.Printf("[Server] Scan error: %v", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "scan_started"})
}
