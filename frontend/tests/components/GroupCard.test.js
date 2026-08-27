import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import GroupCard from '../../src/lib/components/GroupCard.svelte';

describe('GroupCard.svelte', () => {
  const mockGroup = {
    id: 'group-100',
    name: 'Living Room Lights',
    description: 'TV and Ceiling',
    deviceIds: ['wled-1', 'wled-2']
  };

  const mockDevices = [
    { id: 'wled-1', name: 'TV Backlight', isOnline: true },
    { id: 'wled-2', name: 'Ceiling Strip', isOnline: true }
  ];

  it('renders group title, description, and assigned member badges', () => {
    const { getByText } = render(GroupCard, {
      props: {
        group: mockGroup,
        allDevices: mockDevices,
        effects: ['Solid'],
        palettes: ['Default']
      }
    });

    expect(getByText('Living Room Lights')).toBeDefined();
    expect(getByText('TV and Ceiling')).toBeDefined();
    expect(getByText('TV Backlight')).toBeDefined();
    expect(getByText('Ceiling Strip')).toBeDefined();
  });

  it('triggers onTogglePower when power button is clicked', async () => {
    const handleTogglePower = vi.fn();
    const { getByTitle } = render(GroupCard, {
      props: {
        group: mockGroup,
        allDevices: mockDevices,
        effects: ['Solid'],
        palettes: ['Default'],
        onTogglePower: handleTogglePower
      }
    });

    const powerBtn = getByTitle('Toggle Group Power');
    await fireEvent.click(powerBtn);

    expect(handleTogglePower).toHaveBeenCalled();
  });
});
