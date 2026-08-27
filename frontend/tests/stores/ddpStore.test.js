import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getDDPStore } from '../../src/lib/stores/ddpStore.svelte.js';

describe('ddpStore.svelte.js', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('initializes with empty streams map', () => {
    const store = getDDPStore();
    expect(store.isStreaming('device', 'dev-1')).toBe(false);
  });

  it('updates status on WebSocket ddp_status message', () => {
    const store = getDDPStore();
    store.handleWSMessage({
      type: 'ddp_status',
      data: {
        'device:dev-1': { active: true, fps: 60.0, effect: 'rainbow_wave', speed: 1.5, intensity: 0.8 }
      }
    });

    expect(store.isStreaming('device', 'dev-1')).toBe(true);
    expect(store.getStreamStatus('device', 'dev-1').fps).toBe(60.0);
  });

  it('calls REST endpoint on startStream for group target', async () => {
    const store = getDDPStore();
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ active: true, fps: 60.0, effect: 'digital_rain' })
    });

    await store.startStream('group', 'group-desk', 'digital_rain', 1.0, 1.0);

    expect(global.fetch).toHaveBeenCalledWith('/api/v1/ddp/start', expect.any(Object));
    expect(store.isStreaming('group', 'group-desk')).toBe(true);
  });
});
