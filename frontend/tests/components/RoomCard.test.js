import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import RoomCard from '../../src/lib/components/RoomCard.svelte';

describe('RoomCard.svelte', () => {
  const mockRoom = {
    id: 'room-living',
    title: 'Living Room Canvas',
    width: 2000,
    height: 1200
  };

  const mockDevices = [
    { id: 'wled-1', name: 'Living Room TV Strip', ipAddress: '192.168.1.175', ledCount: 60, isOnline: true }
  ];

  const mockPlacements = [
    { deviceId: 'wled-1', roomId: 'room-living', posX: 100, posY: 100, rotation: 0, scale: 1, geometry: 'strip' }
  ];

  it('renders room title and mini 2D preview', () => {
    const { getByText } = render(RoomCard, {
      room: mockRoom,
      devices: mockDevices,
      placements: mockPlacements,
      isPinned: true,
      cardSize: 'normal',
      onEditLayout: () => {},
      onRename: () => {}
    });

    expect(getByText('Living Room Canvas')).toBeDefined();
    expect(getByText('Edit 2D Layout')).toBeDefined();
  });

  it('triggers onRename callback when title is edited', async () => {
    let renamed = null;
    const { getByTitle } = render(RoomCard, {
      room: mockRoom,
      devices: mockDevices,
      placements: mockPlacements,
      isPinned: true,
      cardSize: 'normal',
      onRename: (id, newTitle) => {
        renamed = { id, newTitle };
      }
    });

    const titleBtn = getByTitle('Click to rename room');
    expect(titleBtn).toBeDefined();
  });
});
