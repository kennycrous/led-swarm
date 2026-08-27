import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import DeviceCard from '../../src/lib/components/DeviceCard.svelte';

describe('DeviceCard.svelte', () => {
  const mockDevice = {
    id: 'wled-test-1',
    name: 'Desk LED Strip',
    ipAddress: '192.168.1.175',
    macAddress: 'AA:BB:CC:DD:EE:FF',
    ledCount: 60,
    isOnline: true,
    state: {
      on: true,
      bri: 200,
      seg: [{ id: 0, fx: 0, pal: 0, col: [[6, 182, 212]] }]
    }
  };

  it('renders device title, IP address, and LED count', () => {
    const { getByText } = render(DeviceCard, {
      props: {
        device: mockDevice,
        effects: ['Solid', 'Breathe'],
        palettes: ['Default', 'Colorwaves']
      }
    });

    expect(getByText('Desk LED Strip')).toBeDefined();
    expect(getByText('192.168.1.175 • 60 LEDs')).toBeDefined();
  });

  it('triggers onTogglePower when power button is clicked', async () => {
    const handleTogglePower = vi.fn();
    const { getByTitle } = render(DeviceCard, {
      props: {
        device: mockDevice,
        effects: ['Solid'],
        palettes: ['Default'],
        onTogglePower: handleTogglePower
      }
    });

    const powerBtn = getByTitle('Toggle Power');
    await fireEvent.click(powerBtn);

    expect(handleTogglePower).toHaveBeenCalledWith('wled-test-1');
  });

  it('triggers onSetColor when a color swatch is clicked', async () => {
    const handleSetColor = vi.fn();
    const { getByTitle } = render(DeviceCard, {
      props: {
        device: mockDevice,
        effects: ['Solid'],
        palettes: ['Default'],
        onSetColor: handleSetColor
      }
    });

    const cyanSwatch = getByTitle('Cyan Neon');
    await fireEvent.click(cyanSwatch);

    expect(handleSetColor).toHaveBeenCalledWith('wled-test-1', 6, 182, 212);
  });
});
