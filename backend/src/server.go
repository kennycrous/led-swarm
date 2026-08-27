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
	hub        *Hub
	scanner    *MDNSScanner
	distFS     embed.FS
}

func NewServer(db *Database, wledClient *WLEDClient, hub *Hub, scanner *MDNSScanner, distFS embed.FS) *Server {
	return &Server{
		db:         db,
		wledClient: wledClient,
		hub:        hub,
		scanner:    scanner,
		distFS:     distFS,
	}
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()

	// REST API Endpoints
	mux.HandleFunc("/api/v1/devices", s.handleDevices)
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

	switch r.Method {
	case http.MethodGet:
		devices, err := s.db.GetDevices()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(devices)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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
