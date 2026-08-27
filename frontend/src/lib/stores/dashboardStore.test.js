import { describe, it, expect, beforeEach } from 'vitest';
import { getDashboardStore } from './dashboardStore.svelte.js';

describe('dashboardStore', () => {
  let store;

  beforeEach(() => {
    store = getDashboardStore();
  });

  it('should initialize with empty panels or default state', () => {
    expect(Array.isArray(store.panels)).toBe(true);
    expect(Array.isArray(store.dashboardItems)).toBe(true);
  });

  it('should add a custom panel', async () => {
    await store.addPanel('Living Room Zone');
    const panel = store.panels.find((p) => p.title === 'Living Room Zone');
    expect(panel).toBeDefined();
    expect(panel.title).toBe('Living Room Zone');
  });

  it('should return default normal size for unconfigured items', () => {
    const size = store.getSize('unknown-device-123');
    expect(size).toBe('normal');
  });

  it('should return default empty panelId for unconfigured items', () => {
    const panelId = store.getPanelId('unknown-device-123');
    expect(panelId).toBe('');
  });

  it('should optimistically update size', async () => {
    await store.setSize('wled-device-1', 'wide');
    expect(store.getSize('wled-device-1')).toBe('wide');
  });

  it('should optimistically set panel ID', async () => {
    await store.setPanelId('wled-device-1', 'panel-kitchen');
    expect(store.getPanelId('wled-device-1')).toBe('panel-kitchen');
  });

  it('should delete a custom panel', async () => {
    await store.addPanel('Temporary Panel');
    const added = store.panels.find((p) => p.title === 'Temporary Panel');
    expect(added).toBeDefined();

    await store.deletePanel(added.id);
    const deleted = store.panels.find((p) => p.id === added.id);
    expect(deleted).toBeUndefined();
  });
});
