import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import CreateRoomModal from '../../src/lib/components/CreateRoomModal.svelte';

describe('CreateRoomModal.svelte', () => {
  const mockDevices = [
    { id: 'wled-1', name: 'TV Strip', ipAddress: '192.168.1.101', isOnline: true },
    { id: 'wled-2', name: 'Desk Strip', ipAddress: '192.168.1.102', isOnline: true }
  ];

  it('renders room creation modal when isOpen is true', () => {
    const { getByText, getByPlaceholderText } = render(CreateRoomModal, {
      isOpen: true,
      availableDevices: mockDevices,
      onClose: () => {},
      onCreate: () => {}
    });

    expect(getByText('Create 2D Room Canvas')).toBeDefined();
    expect(getByPlaceholderText('e.g. Living Room, Office, Gaming Den')).toBeDefined();
    expect(getByText('TV Strip')).toBeDefined();
    expect(getByText('Desk Strip')).toBeDefined();
  });

  it('submits form with name, description, and selected strip IDs', async () => {
    const onCreate = vi.fn();
    const { getByPlaceholderText, getByText } = render(CreateRoomModal, {
      isOpen: true,
      availableDevices: mockDevices,
      onClose: () => {},
      onCreate
    });

    const nameInput = getByPlaceholderText('e.g. Living Room, Office, Gaming Den');
    await fireEvent.input(nameInput, { target: { value: 'Living Room Setup' } });

    const descInput = getByPlaceholderText('e.g. Main TV and ceiling strip 2D spatial layout');
    await fireEvent.input(descInput, { target: { value: 'Main TV layout' } });

    // Select TV Strip checkbox
    const stripCheckbox = getByText('TV Strip');
    await fireEvent.click(stripCheckbox);

    const submitBtn = getByText('Create 2D Room');
    await fireEvent.click(submitBtn);

    expect(onCreate).toHaveBeenCalledWith('Living Room Setup', 'Main TV layout', ['wled-1']);
  });
});
