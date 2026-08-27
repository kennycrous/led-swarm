package main

import (
	"log"
	"sync"
)

type CanvasManager struct {
	db         *Database
	hub        *Hub
	placements map[string]CanvasPlacement
	mu         sync.RWMutex
}

func NewCanvasManager(db *Database, hub *Hub) *CanvasManager {
	cm := &CanvasManager{
		db:         db,
		hub:        hub,
		placements: make(map[string]CanvasPlacement),
	}
	_ = cm.loadStoredPlacements()
	return cm
}

func (cm *CanvasManager) loadStoredPlacements() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	stored, err := cm.db.GetCanvasPlacements()
	if err != nil {
		log.Printf("[CanvasManager] Error loading placements from database: %v", err)
		return err
	}

	for _, p := range stored {
		cm.placements[p.DeviceID] = p
	}
	log.Printf("[CanvasManager] Loaded %d canvas placements from database", len(cm.placements))
	return nil
}

func (cm *CanvasManager) GetPlacements() []CanvasPlacement {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	res := make([]CanvasPlacement, 0, len(cm.placements))
	for _, p := range cm.placements {
		res = append(res, p)
	}
	return res
}

func (cm *CanvasManager) SavePlacement(p CanvasPlacement) error {
	cm.mu.Lock()
	if p.Geometry == "" {
		p.Geometry = "strip"
	}
	if p.Scale <= 0 {
		p.Scale = 1.0
	}
	cm.placements[p.DeviceID] = p
	cm.mu.Unlock()

	if err := cm.db.SaveCanvasPlacement(p); err != nil {
		log.Printf("[CanvasManager] Error saving placement for %s: %v", p.DeviceID, err)
		return err
	}

	cm.broadcastCanvasUpdate()
	return nil
}

func (cm *CanvasManager) BatchSavePlacements(placements []CanvasPlacement) error {
	cm.mu.Lock()
	for _, p := range placements {
		if p.Geometry == "" {
			p.Geometry = "strip"
		}
		if p.Scale <= 0 {
			p.Scale = 1.0
		}
		cm.placements[p.DeviceID] = p
		_ = cm.db.SaveCanvasPlacement(p)
	}
	cm.mu.Unlock()

	cm.broadcastCanvasUpdate()
	return nil
}

func (cm *CanvasManager) broadcastCanvasUpdate() {
	if cm.hub != nil {
		cm.hub.BroadcastJSON(map[string]interface{}{
			"type":       "canvas_updated",
			"placements": cm.GetPlacements(),
		})
	}
}
