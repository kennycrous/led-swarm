package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

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
		name = excluded.name,
		ip_address = excluded.ip_address,
		mac_address = excluded.mac_address,
		led_count = excluded.led_count,
		is_online = excluded.is_online,
		last_seen = excluded.last_seen;
	`

	_, err := d.db.Exec(query, dev.ID, dev.Name, dev.IPAddress, dev.MACAddress, dev.LEDCount, dev.IsOnline, time.Now())
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

func (d *Database) Close() error {
	return d.db.Close()
}
