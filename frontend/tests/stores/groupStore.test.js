import { describe, it, expect } from 'vitest';
import { getGroupStore } from '../../src/lib/stores/groupStore.svelte.js';

describe('groupStore', () => {
  it('should initialize empty groups and scenes arrays', () => {
    const store = getGroupStore();
    expect(Array.isArray(store.groups)).toBe(true);
    expect(Array.isArray(store.scenes)).toBe(true);
    expect(store.activeGroupId).toBe('');
  });

  it('should create group optimistically', () => {
    const store = getGroupStore();
    store.createGroup('Office Strip', 'Desk ambient', ['wled-1']);
    expect(store.groups.length).toBeGreaterThan(0);
    const created = store.groups.find((g) => g.name === 'Office Strip');
    expect(created).toBeDefined();
    expect(created.deviceIds).toContain('wled-1');
  });

  it('should delete group optimistically', () => {
    const store = getGroupStore();
    store.createGroup('Temporary Zone', 'To delete', ['wled-99']);
    const target = store.groups.find((g) => g.name === 'Temporary Zone');
    expect(target).toBeDefined();

    store.deleteGroup(target.id);
    const found = store.groups.find((g) => g.id === target.id);
    expect(found).toBeUndefined();
  });
});
