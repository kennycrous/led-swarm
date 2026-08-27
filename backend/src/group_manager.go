package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SceneSnapshot struct {
	DeviceID  string          `json:"deviceId"`
	IPAddress string          `json:"ipAddress"`
	StateJSON json.RawMessage `json:"stateJson"`
}

type GroupManager struct {
	db         *Database
	wledClient *WLEDClient
	devMgr     *DeviceManager
	hub        *Hub
	groups     map[string]*Group
	scenes     map[string]*Scene
	mu         sync.RWMutex
}

func NewGroupManager(db *Database, wledClient *WLEDClient, devMgr *DeviceManager, hub *Hub) *GroupManager {
	gm := &GroupManager{
		db:         db,
		wledClient: wledClient,
		devMgr:     devMgr,
		hub:        hub,
		groups:     make(map[string]*Group),
		scenes:     make(map[string]*Scene),
	}

	if err := gm.loadFromDB(); err != nil {
		log.Printf("[GroupManager] Warning loading groups/scenes: %v", err)
	}

	return gm
}

func (gm *GroupManager) loadFromDB() error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	storedGroups, err := gm.db.GetGroups()
	if err == nil {
		for _, g := range storedGroups {
			groupCopy := g
			gm.groups[g.ID] = &groupCopy
		}
	}

	storedScenes, err := gm.db.GetScenes()
	if err == nil {
		for _, s := range storedScenes {
			sceneCopy := s
			gm.scenes[s.ID] = &sceneCopy
		}
	}

	log.Printf("[GroupManager] Loaded %d groups and %d scenes from database", len(gm.groups), len(gm.scenes))
	return nil
}

func (gm *GroupManager) SaveGroup(id string, name string, description string, deviceIDs []string) (*Group, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if deviceIDs == nil {
		deviceIDs = make([]string, 0)
	}

	g := Group{
		ID:          id,
		Name:        name,
		Description: description,
		DeviceIDs:   deviceIDs,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	if err := gm.db.SaveGroup(g); err != nil {
		return nil, err
	}

	gm.mu.Lock()
	gm.groups[g.ID] = &g
	gm.mu.Unlock()

	gm.broadcastUpdate("group_updated", g)
	return &g, nil
}

func (gm *GroupManager) RenameGroup(id string, newName string) (*Group, error) {
	gm.mu.Lock()
	g, ok := gm.groups[id]
	if !ok {
		gm.mu.Unlock()
		return nil, fmt.Errorf("group not found: %s", id)
	}
	g.Name = newName
	groupCopy := *g
	gm.mu.Unlock()

	if err := gm.db.SaveGroup(groupCopy); err != nil {
		return nil, err
	}

	gm.broadcastUpdate("group_updated", groupCopy)
	return &groupCopy, nil
}

func (gm *GroupManager) DeleteGroup(id string) error {
	if err := gm.db.DeleteGroup(id); err != nil {
		return err
	}

	gm.mu.Lock()
	delete(gm.groups, id)
	gm.mu.Unlock()

	gm.broadcastUpdate("group_deleted", map[string]string{"id": id})
	return nil
}

func (gm *GroupManager) GetGroups() []Group {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	result := make([]Group, 0, len(gm.groups))
	for _, g := range gm.groups {
		result = append(result, *g)
	}
	return result
}

func (gm *GroupManager) GetGroupByID(id string) *Group {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	if g, ok := gm.groups[id]; ok {
		gCopy := *g
		return &gCopy
	}
	return nil
}

// SetGroupState dispatches state payload concurrently across all devices in group
func (gm *GroupManager) SetGroupState(groupID string, rawState json.RawMessage) error {
	gm.mu.RLock()
	group, ok := gm.groups[groupID]
	gm.mu.RUnlock()

	if !ok {
		return fmt.Errorf("group not found: %s", groupID)
	}

	allDevs := gm.devMgr.GetAllDevices()
	devMap := make(map[string]Device)
	for _, d := range allDevs {
		devMap[d.ID] = d
	}

	var wg sync.WaitGroup
	for _, devID := range group.DeviceIDs {
		dev, exists := devMap[devID]
		if !exists || !dev.IsOnline {
			continue
		}

		wg.Add(1)
		go func(targetIP string) {
			defer wg.Done()
			if err := gm.wledClient.SetRawState(targetIP, rawState); err != nil {
				log.Printf("[GroupManager] Error setting state for device %s: %v", targetIP, err)
			}
		}(dev.IPAddress)
	}

	wg.Wait()
	return nil
}

// CaptureScene captures a multi-strip JSON state snapshot across all devices or scoped target
func (gm *GroupManager) CaptureScene(name string, icon string) (*Scene, error) {
	return gm.CaptureScopedScene(name, icon, "global", "")
}

func (gm *GroupManager) CaptureScopedScene(name string, icon string, scopeType string, targetID string) (*Scene, error) {
	devices := gm.devMgr.GetAllDevices()

	if scopeType == "group" && targetID != "" {
		gm.mu.RLock()
		group, ok := gm.groups[targetID]
		gm.mu.RUnlock()
		if ok {
			allowedMap := make(map[string]bool)
			for _, id := range group.DeviceIDs {
				allowedMap[id] = true
			}
			filtered := make([]Device, 0)
			for _, dev := range devices {
				if allowedMap[dev.ID] {
					filtered = append(filtered, dev)
				}
			}
			devices = filtered
		}
	}

	var snapshots []SceneSnapshot
	for _, dev := range devices {
		var rawState json.RawMessage
		if dev.IsOnline {
			fetched, err := gm.wledClient.FetchLiveState(dev.IPAddress)
			if err == nil && len(fetched) > 0 {
				rawState = fetched
			}
		}

		if len(rawState) == 0 {
			stateBytes, err := json.Marshal(dev.State)
			if err == nil {
				rawState = stateBytes
			}
		}

		if len(rawState) > 0 {
			snapshots = append(snapshots, SceneSnapshot{
				DeviceID:  dev.ID,
				IPAddress: dev.IPAddress,
				StateJSON: rawState,
			})
		}
	}

	configBytes, err := json.Marshal(snapshots)
	if err != nil {
		return nil, err
	}

	sceneID := uuid.New().String()
	s := Scene{
		ID:         sceneID,
		Name:       name,
		Icon:       icon,
		ScopeType:  scopeType,
		TargetID:   targetID,
		ConfigJSON: string(configBytes),
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	if err := gm.db.SaveScene(s); err != nil {
		return nil, err
	}

	gm.mu.Lock()
	gm.scenes[s.ID] = &s
	gm.mu.Unlock()

	gm.broadcastUpdate("scene_updated", s)
	return &s, nil
}

// ApplyScene restores multi-strip JSON state snapshot concurrently
func (gm *GroupManager) ApplyScene(sceneID string) error {
	gm.mu.RLock()
	scene, ok := gm.scenes[sceneID]
	gm.mu.RUnlock()

	if !ok {
		return fmt.Errorf("scene not found: %s", sceneID)
	}

	var snapshots []SceneSnapshot
	if err := json.Unmarshal([]byte(scene.ConfigJSON), &snapshots); err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, snap := range snapshots {
		if snap.IPAddress == "" {
			continue
		}
		wg.Add(1)
		go func(targetIP string, raw json.RawMessage) {
			defer wg.Done()
			if err := gm.wledClient.SetRawState(targetIP, raw); err != nil {
				log.Printf("[GroupManager] Error applying scene snapshot to %s: %v", targetIP, err)
			}
		}(snap.IPAddress, snap.StateJSON)
	}

	wg.Wait()
	return nil
}

func (gm *GroupManager) DeleteScene(id string) error {
	if err := gm.db.DeleteScene(id); err != nil {
		return err
	}

	gm.mu.Lock()
	delete(gm.scenes, id)
	gm.mu.Unlock()

	gm.broadcastUpdate("scene_deleted", map[string]string{"id": id})
	return nil
}

func (gm *GroupManager) GetScenes() []Scene {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	result := make([]Scene, 0, len(gm.scenes))
	for _, s := range gm.scenes {
		result = append(result, *s)
	}
	return result
}

func (gm *GroupManager) broadcastUpdate(eventType string, data interface{}) {
	if gm.hub != nil {
		gm.hub.BroadcastJSON(map[string]interface{}{
			"type": eventType,
			"data": data,
		})
	}
}
