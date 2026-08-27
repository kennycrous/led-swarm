package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
)

type Group struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DeviceIDs   []string `json:"deviceIds"`
	CreatedAt   string   `json:"createdAt"`
}

type Scene struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	ConfigJSON string `json:"configJson"`
	CreatedAt  string `json:"createdAt"`
}

type DashboardItem struct {
	ID        string `json:"id"`
	ItemID    string `json:"itemId"`
	ItemType  string `json:"itemType"` // "device", "group", "scene"
	Position  int    `json:"position"`
	Size      string `json:"size"`    // "compact", "normal", "wide"
	PanelID   string `json:"panelId"` // "default" or custom panel id
	IsPinned  bool   `json:"isPinned"`
	CreatedAt string `json:"createdAt"`
}

type DashboardPanel struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Position  int    `json:"position"`
	CreatedAt string `json:"createdAt"`
}

type Database struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite single writer optimization

	d := &Database{db: db}
	if err := d.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return d, nil
}

func (d *Database) initSchema() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	query := `
	CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		ip_address TEXT NOT NULL,
		mac_address TEXT UNIQUE,
		led_count INTEGER DEFAULT 0,
		is_online BOOLEAN DEFAULT FALSE,
		last_seen DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS groups (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS group_devices (
		group_id TEXT NOT NULL,
		device_id TEXT NOT NULL,
		PRIMARY KEY (group_id, device_id),
		FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
		FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS canvas_placements (
		device_id TEXT PRIMARY KEY,
		pos_x REAL NOT NULL DEFAULT 0.0,
		pos_y REAL NOT NULL DEFAULT 0.0,
		rotation REAL NOT NULL DEFAULT 0.0,
		scale_x REAL NOT NULL DEFAULT 1.0,
		scale_y REAL NOT NULL DEFAULT 1.0,
		geometry_type TEXT DEFAULT 'linear',
		FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS scenes (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		icon TEXT,
		config_json TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS dashboard_items (
		id TEXT PRIMARY KEY,
		item_id TEXT UNIQUE NOT NULL,
		item_type TEXT NOT NULL,
		position INTEGER DEFAULT 0,
		size TEXT DEFAULT 'normal',
		panel_id TEXT DEFAULT '',
		is_pinned BOOLEAN DEFAULT TRUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS dashboard_panels (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		position INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("[DB] SQLite database schema initialized successfully")
	return nil
}

func (d *Database) SaveDevice(dev Device) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	query := `
	INSERT INTO devices (id, name, ip_address, mac_address, led_count, is_online, last_seen)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = CASE WHEN excluded.name != '' THEN excluded.name ELSE devices.name END,
		ip_address = excluded.ip_address,
		mac_address = excluded.mac_address,
		led_count = excluded.led_count,
		is_online = excluded.is_online,
		last_seen = excluded.last_seen;
	`

	_, err := d.db.Exec(query, dev.ID, dev.Name, dev.IPAddress, dev.MACAddress, dev.LEDCount, dev.IsOnline, time.Now())
	return err
}

func (d *Database) UpdateDeviceName(id string, name string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec("UPDATE devices SET name = ? WHERE id = ?", name, id)
	return err
}

func (d *Database) DeleteDevice(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec("DELETE FROM devices WHERE id = ?", id)
	return err
}

func (d *Database) GetDevices() ([]Device, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query("SELECT id, name, ip_address, mac_address, led_count, is_online FROM devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var dev Device
		var mac sql.NullString
		if err := rows.Scan(&dev.ID, &dev.Name, &dev.IPAddress, &mac, &dev.LEDCount, &dev.IsOnline); err != nil {
			return nil, err
		}
		if mac.Valid {
			dev.MACAddress = mac.String
		}
		devices = append(devices, dev)
	}

	return devices, nil
}

// Group Database Operations

func (d *Database) SaveGroup(g Group) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO groups (id, name, description)
	VALUES (?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		description = excluded.description;
	`
	if _, err := tx.Exec(query, g.ID, g.Name, g.Description); err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM group_devices WHERE group_id = ?", g.ID); err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO group_devices (group_id, device_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, devID := range g.DeviceIDs {
		if _, err := stmt.Exec(g.ID, devID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) DeleteGroup(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec("DELETE FROM groups WHERE id = ?", id)
	return err
}

func (d *Database) GetGroups() ([]Group, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query("SELECT id, name, description, created_at FROM groups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		var desc sql.NullString
		if err := rows.Scan(&g.ID, &g.Name, &desc, &g.CreatedAt); err != nil {
			return nil, err
		}
		if desc.Valid {
			g.Description = desc.String
		}

		devRows, err := d.db.Query("SELECT device_id FROM group_devices WHERE group_id = ?", g.ID)
		if err == nil {
			g.DeviceIDs = make([]string, 0)
			for devRows.Next() {
				var devID string
				if err := devRows.Scan(&devID); err == nil {
					g.DeviceIDs = append(g.DeviceIDs, devID)
				}
			}
			devRows.Close()
		}

		groups = append(groups, g)
	}

	return groups, nil
}

// Scene Database Operations

func (d *Database) SaveScene(s Scene) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	query := `
	INSERT INTO scenes (id, name, icon, config_json)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		icon = excluded.icon,
		config_json = excluded.config_json;
	`
	_, err := d.db.Exec(query, s.ID, s.Name, s.Icon, s.ConfigJSON)
	return err
}

func (d *Database) DeleteScene(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec("DELETE FROM scenes WHERE id = ?", id)
	return err
}

func (d *Database) GetScenes() ([]Scene, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query("SELECT id, name, icon, config_json, created_at FROM scenes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scenes []Scene
	for rows.Next() {
		var s Scene
		var icon sql.NullString
		if err := rows.Scan(&s.ID, &s.Name, &icon, &s.ConfigJSON, &s.CreatedAt); err != nil {
			return nil, err
		}
		if icon.Valid {
			s.Icon = icon.String
		}
		scenes = append(scenes, s)
	}

	return scenes, nil
}

// Dashboard Items Database Operations

func (d *Database) PinDashboardItem(itemID string, itemType string, isPinned bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	query := `
	INSERT INTO dashboard_items (id, item_id, item_type, is_pinned)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(item_id) DO UPDATE SET
		is_pinned = excluded.is_pinned;
	`
	id := uuid.New().String()
	_, err := d.db.Exec(query, id, itemID, itemType, isPinned)
	return err
}

func (d *Database) UpdateDashboardItemSize(itemID string, size string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	query := `
	INSERT INTO dashboard_items (id, item_id, item_type, size, is_pinned)
	VALUES (?, ?, 'device', ?, TRUE)
	ON CONFLICT(item_id) DO UPDATE SET
		size = excluded.size;
	`
	id := uuid.New().String()
	_, err := d.db.Exec(query, id, itemID, size)
	return err
}

func (d *Database) UpdateDashboardItemPanel(itemID string, panelID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	query := `
	INSERT INTO dashboard_items (id, item_id, item_type, panel_id, is_pinned)
	VALUES (?, ?, 'device', ?, TRUE)
	ON CONFLICT(item_id) DO UPDATE SET
		panel_id = excluded.panel_id;
	`
	id := uuid.New().String()
	_, err := d.db.Exec(query, id, itemID, panelID)
	return err
}

func (d *Database) UpdateDashboardItemPosition(itemID string, position int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec("UPDATE dashboard_items SET position = ? WHERE item_id = ?", position, itemID)
	return err
}

func (d *Database) GetDashboardItems() ([]DashboardItem, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query("SELECT id, item_id, item_type, position, COALESCE(size, 'normal'), COALESCE(panel_id, ''), is_pinned, created_at FROM dashboard_items ORDER BY position ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []DashboardItem
	for rows.Next() {
		var item DashboardItem
		if err := rows.Scan(&item.ID, &item.ItemID, &item.ItemType, &item.Position, &item.Size, &item.PanelID, &item.IsPinned, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (d *Database) SaveDashboardPanel(panel DashboardPanel) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	query := `
	INSERT INTO dashboard_panels (id, title, position)
	VALUES (?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		title = EXCLUDED.title,
		position = EXCLUDED.position
	`
	_, err := d.db.Exec(query, panel.ID, panel.Title, panel.Position)
	return err
}

func (d *Database) GetDashboardPanels() ([]DashboardPanel, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query("SELECT id, title, position, created_at FROM dashboard_panels ORDER BY position ASC, created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var panels []DashboardPanel
	for rows.Next() {
		var p DashboardPanel
		if err := rows.Scan(&p.ID, &p.Title, &p.Position, &p.CreatedAt); err != nil {
			return nil, err
		}
		panels = append(panels, p)
	}
	return panels, nil
}

func (d *Database) DeleteDashboardPanel(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec("DELETE FROM dashboard_panels WHERE id = ?", id)
	if err != nil {
		return err
	}
	_, err = d.db.Exec("UPDATE dashboard_items SET panel_id = '' WHERE panel_id = ?", id)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}
