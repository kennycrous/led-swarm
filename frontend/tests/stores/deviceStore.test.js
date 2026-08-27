import { describe, it, expect } from 'vitest';
import { getDeviceStore } from '../../src/lib/stores/deviceStore.svelte.js';

describe('deviceStore', () => {
  it('should return initial store state arrays', () => {
    const store = getDeviceStore();
    expect(Array.isArray(store.devices)).toBe(true);
    expect(Array.isArray(store.onlineDevices)).toBe(true);
    expect(Array.isArray(store.effects)).toBe(true);
    expect(Array.isArray(store.palettes)).toBe(true);
    expect(store.isScanning).toBe(false);
    expect(store.isConnected).toBe(false);
  });
});
