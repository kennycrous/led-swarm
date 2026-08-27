import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getAudioStore } from '../../src/lib/stores/audioStore.svelte.js';

describe('audioStore.svelte.js', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('initializes with inactive state and empty FFT bins', () => {
    const store = getAudioStore();
    expect(store.isCapturing).toBe(false);
    expect(store.fftBins).toHaveLength(16);
    expect(store.gain).toBe(1.0);
    expect(store.band).toBe('full');
  });

  it('updates gain and band settings', () => {
    const store = getAudioStore();
    store.setGain(2.5);
    store.setBand('bass');

    expect(store.gain).toBe(2.5);
    expect(store.band).toBe('bass');
  });

  it('calculates bass, mid, treble energy from FFT bins', () => {
    const store = getAudioStore();
    // Simulate non-zero FFT energy bins (16 bins)
    store.updateFFTBins([0.8, 0.9, 0.7, 0.6, 0.4, 0.3, 0.2, 0.2, 0.1, 0.1, 0.05, 0.05, 0.02, 0.02, 0.01, 0.01]);

    expect(store.bassEnergy).toBeGreaterThan(0.5);
    expect(store.midEnergy).toBeGreaterThan(0.1);
    expect(store.trebleEnergy).toBeLessThan(0.1);
  });
});
