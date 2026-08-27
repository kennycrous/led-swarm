import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import StripsManagement from '../../src/lib/components/StripsManagement.svelte';

describe('StripsManagement.svelte', () => {
  const mockDevices = [
    {
      id: 'wled-1',
      name: 'Desk Strip',
      ipAddress: '192.168.1.175',
      macAddress: 'AA:BB:CC',
      ledCount: 60,
      isOnline: true
    }
  ];

  it('renders inventory header and registered WLED strips', () => {
    const { getByText } = render(StripsManagement, {
      props: {
        devices: mockDevices,
        isPinned: () => true,
        onTogglePin: () => {},
        onRename: () => {},
        onOpenAddModal: () => {}
      }
    });

    expect(getByText('WLED Light Strips Inventory')).toBeDefined();
    expect(getByText('Desk Strip')).toBeDefined();
    expect(getByText('192.168.1.175')).toBeDefined();
  });
});
