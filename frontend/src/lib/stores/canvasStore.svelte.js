function getWailsApp() {
  if (typeof window !== 'undefined' && window.go && window.go.main && window.go.main.App) {
    return window.go.main.App;
  }
  return null;
}

class CanvasStore {
  placements = $state([]);
  isSaving = $state(false);
  sweepActive = $state(false);

  constructor() {
    this.init();
  }

  async init() {
    await this.fetchPlacements();
  }

  async fetchPlacements() {
    try {
      const wailsApp = getWailsApp();
      if (wailsApp && typeof wailsApp.GetCanvasPlacements === 'function') {
        const data = await wailsApp.GetCanvasPlacements();
        if (data) this.placements = data;
        return;
      }

      const res = await fetch('/api/v1/canvas/placements');
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
    const idx = this.placements.findIndex((p) => p.deviceId === deviceId);
    if (idx !== -1) {
      this.placements[idx] = { ...this.placements[idx], ...updates };
    } else {
      this.placements = [
        ...this.placements,
        {
          deviceId,
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
      const wailsApp = getWailsApp();
      if (wailsApp && typeof wailsApp.BatchSaveCanvasPlacements === 'function') {
        await wailsApp.BatchSaveCanvasPlacements(this.placements);
      } else {
        await fetch('/api/v1/canvas/placements/batch', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.placements)
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
