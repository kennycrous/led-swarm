import { describe, it, expect, beforeEach, vi } from 'vitest';
import { getCanvasStore } from '../../src/lib/stores/canvasStore.svelte.js';

describe('canvasStore', () => {
  beforeEach(() => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockImplementation((url, opts) => {
        if (opts && opts.body) {
          return Promise.resolve({
            ok: true,
            json: () => Promise.resolve({ status: 'success' })
          });
        }
        return Promise.resolve({
          ok: true,
          json: () =>
            Promise.resolve([{ deviceId: 'wled-1', posX: 100, posY: 100, rotation: 0, scale: 1, geometry: 'strip' }])
        });
      })
    );
  });

  it('should initialize empty placements array', () => {
    const store = getCanvasStore();
    expect(Array.isArray(store.placements)).toBe(true);
  });

  it('should update placement position optimistically', () => {
    const store = getCanvasStore();
    store.currentRoomId = 'test-room';
    store.placements = [
      { deviceId: 'wled-1', roomId: 'test-room', posX: 100, posY: 100, rotation: 0, scale: 1, geometry: 'strip' }
    ];

    store.updatePlacement('wled-1', { posX: 250, posY: 350 });

    const p = store.placements.find((item) => item.deviceId === 'wled-1');
    expect(p).toBeDefined();
    expect(p.posX).toBe(250);
    expect(p.posY).toBe(350);
  });
});
