package main

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/google/uuid"
)

type CanvasManager struct {
	db         *Database
	hub        *Hub
	rooms      map[string]CanvasRoom
	placements map[string]CanvasPlacement // Key: "roomId:deviceId"
	mu         sync.RWMutex
}

func NewCanvasManager(db *Database, hub *Hub) *CanvasManager {
	cm := &CanvasManager{
		db:         db,
		hub:        hub,
		rooms:      make(map[string]CanvasRoom),
		placements: make(map[string]CanvasPlacement),
	}
	_ = cm.loadStoredData()
	return cm
}

func (cm *CanvasManager) loadStoredData() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	storedRooms, err := cm.db.GetCanvasRooms()
	if err == nil {
		for _, r := range storedRooms {
			cm.rooms[r.ID] = r
		}
	}

	storedPlacements, err := cm.db.GetCanvasPlacements()
	if err == nil {
		for _, p := range storedPlacements {
			if p.RoomID == "" {
				p.RoomID = "default"
			}
			key := p.RoomID + ":" + p.DeviceID
			cm.placements[key] = p
		}
	}

	log.Printf("[CanvasManager] Loaded %d canvas rooms and %d placements from database", len(cm.rooms), len(cm.placements))
	return nil
}

func (cm *CanvasManager) CreateRoom(title string, description string, width int, height int, deviceIDs []string) (CanvasRoom, error) {
	cm.mu.Lock()
	if width <= 0 {
		width = 2000
	}
	if height <= 0 {
		height = 1200
	}
	room := CanvasRoom{
		ID:          "room-" + uuid.New().String()[:8],
		Title:       title,
		Description: description,
		Width:       width,
		Height:      height,
	}
	cm.rooms[room.ID] = room

	// Create initial placements for assigned devices
	for i, devID := range deviceIDs {
		p := CanvasPlacement{
			DeviceID: devID,
			RoomID:   room.ID,
			PosX:     float64(150 + (i%4)*250),
			PosY:     float64(150 + (i/4)*180),
			Rotation: 0,
			Scale:    1.0,
			Geometry: "strip",
		}
		key := room.ID + ":" + devID
		cm.placements[key] = p
		_ = cm.db.SaveCanvasPlacement(p)
	}
	cm.mu.Unlock()

	if err := cm.db.SaveCanvasRoom(room); err != nil {
		return CanvasRoom{}, err
	}

	cm.broadcastCanvasUpdate()
	return room, nil
}

func (cm *CanvasManager) UpdateRoomTitle(id string, newTitle string) (*CanvasRoom, error) {
	cm.mu.Lock()
	room, ok := cm.rooms[id]
	if !ok {
		cm.mu.Unlock()
		return nil, fmt.Errorf("room not found: %s", id)
	}
	room.Title = newTitle
	cm.rooms[id] = room
	cm.mu.Unlock()

	if err := cm.db.SaveCanvasRoom(room); err != nil {
		return nil, err
	}

	cm.broadcastCanvasUpdate()
	return &room, nil
}

func (cm *CanvasManager) GetRooms() []CanvasRoom {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	res := make([]CanvasRoom, 0, len(cm.rooms))
	for _, r := range cm.rooms {
		res = append(res, r)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt < res[j].CreatedAt || (res[i].CreatedAt == res[j].CreatedAt && res[i].ID < res[j].ID)
	})
	return res
}

func (cm *CanvasManager) DeleteRoom(id string) error {
	cm.mu.Lock()
	delete(cm.rooms, id)
	for key, p := range cm.placements {
		if p.RoomID == id {
			delete(cm.placements, key)
		}
	}
	cm.mu.Unlock()

	if err := cm.db.DeleteCanvasRoom(id); err != nil {
		return err
	}

	cm.broadcastCanvasUpdate()
	return nil
}

func (cm *CanvasManager) GetPlacementsForRoom(roomID string) []CanvasPlacement {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	res := make([]CanvasPlacement, 0)
	if roomID == "" || roomID == "all" || roomID == "*" {
		for _, p := range cm.placements {
			res = append(res, p)
		}
		return res
	}

	for _, p := range cm.placements {
		if p.RoomID == roomID {
			res = append(res, p)
		}
	}
	return res
}

func (cm *CanvasManager) SavePlacement(p CanvasPlacement) error {
	cm.mu.Lock()
	if p.RoomID == "" {
		p.RoomID = "default"
	}
	if p.Geometry == "" {
		p.Geometry = "strip"
	}
	if p.Scale <= 0 {
		p.Scale = 1.0
	}
	key := p.RoomID + ":" + p.DeviceID
	cm.placements[key] = p
	cm.mu.Unlock()

	if err := cm.db.SaveCanvasPlacement(p); err != nil {
		log.Printf("[CanvasManager] Error saving placement for %s in room %s: %v", p.DeviceID, p.RoomID, err)
		return err
	}

	cm.broadcastCanvasUpdate()
	return nil
}

func (cm *CanvasManager) BatchSavePlacements(roomID string, placements []CanvasPlacement) error {
	if roomID == "" {
		roomID = "default"
	}
	cm.mu.Lock()
	for _, p := range placements {
		p.RoomID = roomID
		if p.Geometry == "" {
			p.Geometry = "strip"
		}
		if p.Scale <= 0 {
			p.Scale = 1.0
		}
		key := roomID + ":" + p.DeviceID
		cm.placements[key] = p
		_ = cm.db.SaveCanvasPlacement(p)
	}
	cm.mu.Unlock()

	cm.broadcastCanvasUpdate()
	return nil
}

func (cm *CanvasManager) broadcastCanvasUpdate() {
	if cm.hub != nil {
		cm.hub.BroadcastJSON(map[string]interface{}{
			"type":  "canvas_updated",
			"rooms": cm.GetRooms(),
		})
	}
}
