function getWailsApp() {
  if (typeof window !== 'undefined' && window.go && window.go.main && window.go.main.App) {
    return window.go.main.App;
  }
  return null;
}

function normalizeRoom(r) {
  if (!r) return null;
  return {
    id: r.id || r.ID || '',
    title: r.title || r.Title || 'Room',
    description: r.description || r.Description || '',
    width: r.width || r.Width || 2000,
    height: r.height || r.Height || 1200,
    createdAt: r.createdAt || r.CreatedAt || ''
  };
}

function normalizePlacement(p) {
  if (!p) return null;
  return {
    deviceId: p.deviceId || p.DeviceID || '',
    roomId: p.roomId || p.RoomID || '',
    posX: p.posX ?? p.PosX ?? 100,
    posY: p.posY ?? p.PosY ?? 100,
    rotation: p.rotation ?? p.Rotation ?? 0,
    scale: p.scale ?? p.Scale ?? 1.0,
    geometry: p.geometry || p.Geometry || 'strip'
  };
}

class CanvasStore {
  rooms = $state([]);
  currentRoomId = $state(null);
  placements = $state([]);
  isSaving = $state(false);
  sweepActive = $state(false);

  constructor() {
    this.init();
  }

  async init() {
    await this.fetchRooms();
    if (this.rooms.length > 0) {
      this.currentRoomId = this.rooms[0].id;
    } else {
      this.currentRoomId = null;
    }
    await this.fetchPlacements();
  }

  async fetchRooms() {
    try {
      const wailsApp = getWailsApp();
      let rawData = null;
      if (wailsApp && typeof wailsApp.GetCanvasRooms === 'function') {
        rawData = await wailsApp.GetCanvasRooms();
      } else {
        const res = await fetch('/api/v1/canvas/rooms');
        if (res.ok) rawData = await res.json();
      }

      if (Array.isArray(rawData)) {
        this.rooms = rawData.map(normalizeRoom).filter(Boolean);
        if (this.rooms.length > 0 && !this.rooms.some((r) => r.id === this.currentRoomId)) {
          this.currentRoomId = this.rooms[0].id;
        } else if (this.rooms.length === 0) {
          this.currentRoomId = null;
        }
      }
    } catch (err) {
      console.warn('[CanvasStore] Error fetching rooms:', err);
    }
  }

  async createRoom(title, description = '', deviceIds = [], width = 2000, height = 1200) {
    try {
      const wailsApp = getWailsApp();
      let rawRoom = null;
      if (wailsApp && typeof wailsApp.CreateCanvasRoom === 'function') {
        rawRoom = await wailsApp.CreateCanvasRoom(title, description, width, height, deviceIds);
      } else {
        const res = await fetch('/api/v1/canvas/rooms', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ title, description, width, height, deviceIds })
        });
        if (res.ok) rawRoom = await res.json();
      }

      const newRoom = normalizeRoom(rawRoom);
      if (newRoom) {
        this.rooms = [...this.rooms.filter((r) => r.id !== newRoom.id), newRoom];
        this.currentRoomId = newRoom.id;
        await this.fetchPlacements();
      }
      return newRoom;
    } catch (err) {
      console.error('[CanvasStore] Error creating room:', err);
    }
  }

  async deleteRoom(id) {
    try {
      const wailsApp = getWailsApp();
      if (wailsApp && typeof wailsApp.DeleteCanvasRoom === 'function') {
        await wailsApp.DeleteCanvasRoom(id);
      } else {
        await fetch(`/api/v1/canvas/rooms/${id}`, { method: 'DELETE' });
      }

      this.rooms = this.rooms.filter((r) => r.id !== id);
      this.placements = this.placements.filter((p) => p.roomId !== id);
      if (this.currentRoomId === id) {
        this.currentRoomId = this.rooms[0]?.id || null;
      }
    } catch (err) {
      console.error('[CanvasStore] Error deleting room:', err);
    }
  }

  async selectRoom(roomId) {
    this.currentRoomId = roomId;
    await this.fetchPlacements(roomId);
  }

  async renameRoom(roomId, newTitle) {
    try {
      const found = this.rooms.find((r) => r.id === roomId);
      if (found) found.title = newTitle;

      const wailsApp = getWailsApp();
      if (wailsApp && typeof wailsApp.UpdateRoomTitle === 'function') {
        await wailsApp.UpdateRoomTitle(roomId, newTitle);
      } else {
        await fetch('/api/v1/canvas/rooms/rename', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ id: roomId, title: newTitle })
        });
      }
    } catch (err) {
      console.error('[CanvasStore] Error renaming room:', err);
    }
  }

  async fetchPlacements(roomId = '') {
    try {
      const wailsApp = getWailsApp();
      let rawData = null;
      if (wailsApp && typeof wailsApp.GetCanvasPlacements === 'function') {
        rawData = await wailsApp.GetCanvasPlacements(roomId || '');
      } else {
        const url = roomId ? `/api/v1/canvas/placements?roomId=${roomId}` : '/api/v1/canvas/placements';
        const res = await fetch(url);
        if (res.ok) rawData = await res.json();
      }

      if (Array.isArray(rawData)) {
        const normalized = rawData.map(normalizePlacement).filter(Boolean);
        if (!roomId) {
          this.placements = normalized;
        } else {
          const otherPlacements = this.placements.filter((p) => p.roomId !== roomId);
          this.placements = [...otherPlacements, ...normalized];
        }
      }
    } catch (err) {
      console.warn('[CanvasStore] Error fetching placements:', err);
    }
  }

  updatePlacement(deviceId, updates) {
    if (!this.currentRoomId) return;
    const idx = this.placements.findIndex((p) => p.deviceId === deviceId && p.roomId === this.currentRoomId);
    if (idx !== -1) {
      this.placements[idx] = { ...this.placements[idx], ...updates, roomId: this.currentRoomId };
    } else {
      this.placements = [
        ...this.placements,
        {
          deviceId,
          roomId: this.currentRoomId,
          posX: updates.posX ?? 100,
          posY: updates.posY ?? 100,
          rotation: updates.rotation ?? 0,
          scale: updates.scale ?? 1.0,
          geometry: updates.geometry ?? 'strip'
        }
      ];
    }
  }

  async savePlacements() {
    if (!this.currentRoomId) return;
    this.isSaving = true;
    try {
      const roomPlacements = this.placements
        .filter((p) => p.roomId === this.currentRoomId)
        .map((p) => ({ ...p, roomId: this.currentRoomId }));

      const payload = {
        roomId: this.currentRoomId,
        placements: roomPlacements
      };

      const wailsApp = getWailsApp();
      if (wailsApp && typeof wailsApp.BatchSaveCanvasPlacements === 'function') {
        await wailsApp.BatchSaveCanvasPlacements(this.currentRoomId, payload.placements);
      } else {
        await fetch('/api/v1/canvas/placements/batch', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
      }
    } catch (err) {
      console.error('[CanvasStore] Error saving placements:', err);
    } finally {
      this.isSaving = false;
    }
  }

  triggerSpatialSweep() {
    this.sweepActive = true;
    setTimeout(() => {
      this.sweepActive = false;
    }, 3000);
  }
}

let instance = null;
export function getCanvasStore() {
  if (!instance) {
    instance = new CanvasStore();
  }
  return instance;
}
