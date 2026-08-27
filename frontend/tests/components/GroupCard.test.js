import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import GroupCard from '../../src/lib/components/GroupCard.svelte';

describe('GroupCard.svelte', () => {
  const mockGroup = {
    id: 'group-1',
    name: 'Living Room Lights',
    description: 'TV and Ceiling',
    deviceIds: ['wled-1', 'wled-2']
  };

  const mockDevices = [
    { id: 'wled-1', name: 'TV Backlight', isOnline: true },
    { id: 'wled-2', name: 'Ceiling Strip', isOnline: true }
  ];

  it('renders group card details and assigned strips', () => {
    const { getByText } = render(GroupCard, {
      group: mockGroup,
      allDevices: mockDevices
    });

    expect(getByText('Living Room Lights')).toBeDefined();
    expect(getByText('TV Backlight')).toBeDefined();
    expect(getByText('Ceiling Strip')).toBeDefined();
  });

  it('triggers onTogglePower when power button is clicked', async () => {
    const onTogglePower = vi.fn();
    const { getByTitle } = render(GroupCard, {
      group: mockGroup,
      allDevices: mockDevices,
      onTogglePower
    });

    const powerBtn = getByTitle('Toggle Group Power');
    await fireEvent.click(powerBtn);

    expect(onTogglePower).toHaveBeenCalledWith('group-1', false);
  });
});
