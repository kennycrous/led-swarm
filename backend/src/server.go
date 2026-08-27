package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

// Server handles headless HTTP web API and static file serving
type Server struct {
	db           *Database
	wledClient   *WLEDClient
	devMgr       *DeviceManager
	groupMgr     *GroupManager
	dashboardMgr *DashboardManager
	canvasMgr    *CanvasManager
	hub          *Hub
	scanner      *MDNSScanner
	distFS       embed.FS
}

func NewServer(db *Database, wledClient *WLEDClient, devMgr *DeviceManager, groupMgr *GroupManager, dashboardMgr *DashboardManager, canvasMgr *CanvasManager, hub *Hub, scanner *MDNSScanner, distFS embed.FS) *Server {
	return &Server{
		db:           db,
		wledClient:   wledClient,
		devMgr:       devMgr,
		groupMgr:     groupMgr,
		dashboardMgr: dashboardMgr,
		canvasMgr:    canvasMgr,
		hub:          hub,
		scanner:      scanner,
		distFS:       distFS,
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

	// Group & Scene API Endpoints
	mux.HandleFunc("/api/v1/groups", s.handleGroups)
	mux.HandleFunc("/api/v1/groups/", s.handleGroups)
	mux.HandleFunc("/api/v1/groups/state", s.handleSetGroupState)
	mux.HandleFunc("/api/v1/scenes", s.handleScenes)
	mux.HandleFunc("/api/v1/scenes/", s.handleScenes)
	mux.HandleFunc("/api/v1/scenes/capture", s.handleCaptureScene)
	mux.HandleFunc("/api/v1/scenes/apply", s.handleApplyScene)

	// Dashboard API Endpoints
	mux.HandleFunc("/api/v1/dashboard/items", s.handleDashboardItems)
	mux.HandleFunc("/api/v1/dashboard/pin", s.handlePinDashboardItem)
	mux.HandleFunc("/api/v1/dashboard/size", s.handleSetDashboardItemSize)
	mux.HandleFunc("/api/v1/dashboard/panel", s.handleSetDashboardItemPanel)
	mux.HandleFunc("/api/v1/dashboard/panels", s.handleDashboardPanels)
	mux.HandleFunc("/api/v1/dashboard/panels/", s.handleDashboardPanels)
	mux.HandleFunc("/api/v1/dashboard/reorder", s.handleReorderDashboardItems)

	// 2D Canvas API Endpoints
	mux.HandleFunc("/api/v1/canvas/placements", s.handleCanvasPlacements)
	mux.HandleFunc("/api/v1/canvas/placement", s.handleSaveCanvasPlacement)
	mux.HandleFunc("/api/v1/canvas/placements/batch", s.handleBatchSaveCanvasPlacements)

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

// Group Handlers

func (s *Server) handleGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		groups := s.groupMgr.GetGroups()
		json.NewEncoder(w).Encode(groups)

	case http.MethodPost:
		var req struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			DeviceIDs   []string `json:"deviceIds"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, "Invalid parameters", http.StatusBadRequest)
			return
		}
		g, err := s.groupMgr.SaveGroup(req.ID, req.Name, req.Description, req.DeviceIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(g)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing group id parameter", http.StatusBadRequest)
			return
		}
		if err := s.groupMgr.DeleteGroup(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "group_deleted"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleSetGroupState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		GroupID string          `json:"groupId"`
		State   json.RawMessage `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.GroupID == "" {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}
	if err := s.groupMgr.SetGroupState(req.GroupID, req.State); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "group_state_updated"})
}

// Scene Handlers

func (s *Server) handleScenes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		scenes := s.groupMgr.GetScenes()
		json.NewEncoder(w).Encode(scenes)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing scene id parameter", http.StatusBadRequest)
			return
		}
		if err := s.groupMgr.DeleteScene(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "scene_deleted"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleCaptureScene(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Invalid scene name provided", http.StatusBadRequest)
		return
	}
	sc, err := s.groupMgr.CaptureScene(req.Name, req.Icon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sc)
}

func (s *Server) handleApplyScene(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" {
		http.Error(w, "Invalid scene ID", http.StatusBadRequest)
		return
	}
	if err := s.groupMgr.ApplyScene(req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "scene_applied"})
}

// Dashboard Handlers

func (s *Server) handleDashboardItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	items := s.dashboardMgr.GetItems()
	json.NewEncoder(w).Encode(items)
}

func (s *Server) handlePinDashboardItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ItemID   string `json:"itemId"`
		ItemType string `json:"itemType"`
		IsPinned bool   `json:"isPinned"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ItemID == "" {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}
	item, err := s.dashboardMgr.PinItem(req.ItemID, req.ItemType, req.IsPinned)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func (s *Server) handleSetDashboardItemSize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ItemID string `json:"itemId"`
		Size   string `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ItemID == "" {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}
	item, err := s.dashboardMgr.SetItemSize(req.ItemID, req.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func (s *Server) handleSetDashboardItemPanel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ItemID  string `json:"itemId"`
		PanelID string `json:"panelId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ItemID == "" {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}
	item, err := s.dashboardMgr.SetItemPanel(req.ItemID, req.PanelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func (s *Server) handleReorderDashboardItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ItemIDs []string `json:"itemIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}
	if err := s.dashboardMgr.ReorderItems(req.ItemIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "reordered"})
}

func (s *Server) handleDashboardPanels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		panels, err := s.dashboardMgr.GetPanels()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(panels)
	case http.MethodPost:
		var req struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
			http.Error(w, "Invalid panel title", http.StatusBadRequest)
			return
		}
		panel, err := s.dashboardMgr.AddPanel(req.ID, req.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(panel)
	case http.MethodDelete:
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/dashboard/panels/")
		if id == "" {
			var req struct {
				ID string `json:"id"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			id = req.ID
		}
		if id == "" {
			http.Error(w, "Missing panel ID", http.StatusBadRequest)
			return
		}
		if err := s.dashboardMgr.DeletePanel(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleCanvasPlacements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.canvasMgr.GetPlacements())
}

func (s *Server) handleSaveCanvasPlacement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var placement CanvasPlacement
	if err := json.NewDecoder(r.Body).Decode(&placement); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	if err := s.canvasMgr.SavePlacement(placement); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleBatchSaveCanvasPlacements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var placements []CanvasPlacement
	if err := json.NewDecoder(r.Body).Decode(&placements); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	if err := s.canvasMgr.BatchSavePlacements(placements); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
