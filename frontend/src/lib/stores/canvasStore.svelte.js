function getWailsApp() {
  if (typeof window !== 'undefined' && window.go && window.go.main && window.go.main.App) {
    return window.go.main.App;
  }
  return null;
}

class CanvasStore {
  rooms = $state([{ id: 'default', title: 'Main Room Canvas', width: 2000, height: 1200 }]);
  currentRoomId = $state('default');
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
    }
    await this.fetchPlacements();
  }

  async fetchRooms() {
    try {
      const wailsApp = getWailsApp();
      if (wailsApp && typeof wailsApp.GetCanvasRooms === 'function') {
        const data = await wailsApp.GetCanvasRooms();
        if (data && data.length > 0) this.rooms = data;
        return;
      }

      const res = await fetch('/api/v1/canvas/rooms');
      if (res.ok) {
        const data = await res.json();
        if (Array.isArray(data) && data.length > 0) {
          this.rooms = data;
        }
      }
    } catch (err) {
      console.warn('[CanvasStore] Error fetching rooms:', err);
    }
  }

  async createRoom(title, description = '', deviceIds = [], width = 2000, height = 1200) {
    try {
      const wailsApp = getWailsApp();
      let newRoom = null;
      if (wailsApp && typeof wailsApp.CreateCanvasRoom === 'function') {
        newRoom = await wailsApp.CreateCanvasRoom(title, description, width, height, deviceIds);
      } else {
        const res = await fetch('/api/v1/canvas/rooms', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ title, description, width, height, deviceIds })
        });
        if (res.ok) newRoom = await res.json();
      }

      if (newRoom) {
        this.rooms = [...this.rooms, newRoom];
        this.currentRoomId = newRoom.id;
        await this.fetchPlacements(newRoom.id);
      }
      return newRoom;
    } catch (err) {
      console.error('[CanvasStore] Error creating room:', err);
    }
  }

  async deleteRoom(id) {
    if (id === 'default') return;
    try {
      const wailsApp = getWailsApp();
      if (wailsApp && typeof wailsApp.DeleteCanvasRoom === 'function') {
        await wailsApp.DeleteCanvasRoom(id);
      } else {
        await fetch(`/api/v1/canvas/rooms/${id}`, { method: 'DELETE' });
      }

      this.rooms = this.rooms.filter((r) => r.id !== id);
      if (this.currentRoomId === id) {
        this.currentRoomId = this.rooms[0]?.id || 'default';
        await this.fetchPlacements(this.currentRoomId);
      }
    } catch (err) {
      console.error('[CanvasStore] Error deleting room:', err);
    }
  }

  async selectRoom(roomId) {
    this.currentRoomId = roomId;
    await this.fetchPlacements(roomId);
  }

  async fetchPlacements(roomId = '') {
    try {
      const wailsApp = getWailsApp();
      if (wailsApp && typeof wailsApp.GetCanvasPlacements === 'function') {
        const data = await wailsApp.GetCanvasPlacements(roomId || '');
        if (data) this.placements = data;
        return;
      }

      const url = roomId ? `/api/v1/canvas/placements?roomId=${roomId}` : '/api/v1/canvas/placements';
      const res = await fetch(url);
      if (res.ok) {
        const data = await res.json();
        if (Array.isArray(data)) {
          this.placements = data;
        }
      }
    } catch (err) {
      console.warn('[CanvasStore] Error fetching placements:', err);
    }
  }

  updatePlacement(deviceId, updates) {
    const idx = this.placements.findIndex(
      (p) => p.deviceId === deviceId && (p.roomId === this.currentRoomId || !p.roomId)
    );
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
    this.isSaving = true;
    try {
      const payload = {
        roomId: this.currentRoomId,
        placements: this.placements.map((p) => ({ ...p, roomId: this.currentRoomId }))
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
