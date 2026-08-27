import { describe, it, expect, vi } from 'vitest';
import { render } from '@testing-library/svelte';
import CanvasEditor from '../../src/lib/components/CanvasEditor.svelte';

describe('CanvasEditor.svelte', () => {
  const mockDevices = [{ id: 'wled-1', name: 'Desk Strip', ipAddress: '192.168.1.175', ledCount: 60, isOnline: true }];

  const mockPlacements = [{ deviceId: 'wled-1', posX: 120, posY: 180, rotation: 0, scale: 1, geometry: 'strip' }];

  it('renders canvas editor header and room grid', () => {
    const { getByText } = render(CanvasEditor, {
      devices: mockDevices,
      placements: mockPlacements,
      onSavePlacements: () => {}
    });

    expect(getByText('2D Room Layout Canvas')).toBeDefined();
    expect(getByText('Desk Strip')).toBeDefined();
  });
});
