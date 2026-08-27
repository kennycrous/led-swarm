import { describe, it, expect, beforeEach, vi } from 'vitest';
import { getGroupStore } from '../../src/lib/stores/groupStore.svelte.js';

describe('groupStore', () => {
  beforeEach(() => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockImplementation((url, opts) => {
        if (opts && opts.body) {
          const payload = JSON.parse(opts.body);
          return Promise.resolve({
            ok: true,
            json: () =>
              Promise.resolve({
                id: `group-${Math.random()}`,
                name: payload.name,
                description: payload.description,
                deviceIds: payload.deviceIds
              })
          });
        }
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([])
        });
      })
    );
  });

  it('should initialize empty groups and scenes arrays', () => {
    const store = getGroupStore();
    expect(Array.isArray(store.groups)).toBe(true);
    expect(Array.isArray(store.scenes)).toBe(true);
  });

  it('should create group optimistically', async () => {
    const store = getGroupStore();
    await store.createGroup('Office Strip', 'Desk ambient', ['wled-1']);
    expect(store.groups.length).toBeGreaterThan(0);
    const created = store.groups.find((g) => g.name === 'Office Strip');
    expect(created).toBeDefined();
    expect(created.deviceIds).toContain('wled-1');
  });

  it('should delete group optimistically', async () => {
    const store = getGroupStore();
    await store.createGroup('Temporary Zone', 'To delete', ['wled-99']);
    const target = store.groups.find((g) => g.name === 'Temporary Zone');
    expect(target).toBeDefined();

    await store.deleteGroup(target.id);
    const found = store.groups.find((g) => g.id === target.id);
    expect(found).toBeUndefined();
  });
});
