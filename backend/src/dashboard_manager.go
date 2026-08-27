package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type DashboardManager struct {
	db    *Database
	hub   *Hub
	items map[string]*DashboardItem // Keyed by itemID
	mu    sync.RWMutex
}

func NewDashboardManager(db *Database, hub *Hub) *DashboardManager {
	dm := &DashboardManager{
		db:    db,
		hub:   hub,
		items: make(map[string]*DashboardItem),
	}

	if err := dm.loadFromDB(); err != nil {
		log.Printf("[DashboardManager] Warning loading dashboard items: %v", err)
	}

	return dm
}

func (dm *DashboardManager) loadFromDB() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	storedItems, err := dm.db.GetDashboardItems()
	if err == nil {
		for _, item := range storedItems {
			copyItem := item
			dm.items[item.ItemID] = &copyItem
		}
	}

	log.Printf("[DashboardManager] Loaded %d dashboard item rules from database", len(dm.items))
	return nil
}

func (dm *DashboardManager) PinItem(itemID string, itemType string, isPinned bool) (*DashboardItem, error) {
	dm.mu.Lock()
	existingItem, ok := dm.items[itemID]
	panelID := ""
	size := "normal"
	if ok {
		panelID = existingItem.PanelID
		size = existingItem.Size
	}
	dm.mu.Unlock()

	if err := dm.db.PinDashboardItem(itemID, itemType, isPinned); err != nil {
		return nil, err
	}

	dm.mu.Lock()
	item, ok := dm.items[itemID]
	if !ok {
		item = &DashboardItem{
			ItemID:   itemID,
			ItemType: itemType,
			IsPinned: isPinned,
			PanelID:  panelID,
			Size:     size,
		}
		dm.items[itemID] = item
	} else {
		item.IsPinned = isPinned
	}
	itemCopy := *item
	dm.mu.Unlock()

	dm.broadcastUpdate("dashboard_pin_updated", itemCopy)
	return &itemCopy, nil
}

func (dm *DashboardManager) SetItemSize(itemID string, size string) (*DashboardItem, error) {
	if size == "" {
		size = "normal"
	}
	if err := dm.db.UpdateDashboardItemSize(itemID, size); err != nil {
		return nil, err
	}

	dm.mu.Lock()
	item, ok := dm.items[itemID]
	if !ok {
		item = &DashboardItem{
			ItemID:   itemID,
			ItemType: "device",
			Size:     size,
			IsPinned: true,
		}
		dm.items[itemID] = item
	} else {
		item.Size = size
	}
	itemCopy := *item
	dm.mu.Unlock()

	dm.broadcastUpdate("dashboard_size_updated", itemCopy)
	return &itemCopy, nil
}

func (dm *DashboardManager) SetItemPanel(itemID string, panelID string) (*DashboardItem, error) {
	if err := dm.db.UpdateDashboardItemPanel(itemID, panelID); err != nil {
		return nil, err
	}

	dm.mu.Lock()
	item, ok := dm.items[itemID]
	if !ok {
		item = &DashboardItem{
			ItemID:   itemID,
			ItemType: "device",
			PanelID:  panelID,
			IsPinned: true,
			Size:     "normal",
		}
		dm.items[itemID] = item
	} else {
		item.PanelID = panelID
	}
	itemCopy := *item
	dm.mu.Unlock()

	dm.broadcastUpdate("dashboard_panel_updated", itemCopy)
	return &itemCopy, nil
}

func (dm *DashboardManager) ReorderItems(itemIDs []string) error {
	dm.mu.Lock()
	for pos, itemID := range itemIDs {
		_ = dm.db.UpdateDashboardItemPosition(itemID, pos)
		if item, ok := dm.items[itemID]; ok {
			item.Position = pos
		}
	}
	dm.mu.Unlock()

	dm.broadcastUpdate("dashboard_reordered", itemIDs)
	return nil
}

func (dm *DashboardManager) GetItems() []DashboardItem {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	result := make([]DashboardItem, 0, len(dm.items))
	for _, item := range dm.items {
		result = append(result, *item)
	}
	return result
}

func (dm *DashboardManager) AutoPinIfNew(itemID string, itemType string) {
	dm.mu.RLock()
	_, exists := dm.items[itemID]
	dm.mu.RUnlock()

	if !exists {
		dm.PinItem(itemID, itemType, true)
	}
}

func (dm *DashboardManager) AddPanel(id string, title string) (*DashboardPanel, error) {
	if id == "" {
		id = fmt.Sprintf("panel-%d", time.Now().UnixNano())
	}
	panel := DashboardPanel{
		ID:    id,
		Title: title,
	}
	if err := dm.db.SaveDashboardPanel(panel); err != nil {
		return nil, err
	}
	dm.broadcastUpdate("dashboard_panel_added", panel)
	return &panel, nil
}

func (dm *DashboardManager) GetPanels() ([]DashboardPanel, error) {
	return dm.db.GetDashboardPanels()
}

func (dm *DashboardManager) DeletePanel(id string) error {
	if err := dm.db.DeleteDashboardPanel(id); err != nil {
		return err
	}
	dm.broadcastUpdate("dashboard_panel_deleted", map[string]string{"id": id})
	return nil
}

func (dm *DashboardManager) broadcastUpdate(eventType string, data interface{}) {
	if dm.hub != nil {
		dm.hub.BroadcastJSON(map[string]interface{}{
			"type": eventType,
			"data": data,
		})
	}
}
