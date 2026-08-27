import { describe, it, expect, vi } from 'vitest';
import { render } from '@testing-library/svelte';
import CanvasEditor from '../../src/lib/components/CanvasEditor.svelte';

describe('CanvasEditor.svelte', () => {
  const mockDevices = [{ id: 'wled-1', name: 'Desk Strip', ipAddress: '192.168.1.175', ledCount: 60, isOnline: true }];

  const mockPlacements = [{ deviceId: 'wled-1', posX: 120, posY: 180, rotation: 0, scale: 1, geometry: 'strip' }];

  it('renders canvas editor header and room grid', () => {
    const { getByText } = render(CanvasEditor, {
      devices: mockDevices,
      placements: mockPlacements.map((p) => ({ ...p, roomId: 'room-1' })),
      rooms: [{ id: 'room-1', title: 'Living Room', width: 2000, height: 1200 }],
      currentRoomId: 'room-1',
      onSavePlacements: () => {}
    });

    expect(getByText('Living Room 2D Layout Canvas')).toBeDefined();
    expect(getByText('Desk Strip')).toBeDefined();
  });
});
