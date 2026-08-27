import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import CreateGroupModal from '../../src/lib/components/CreateGroupModal.svelte';

describe('CreateGroupModal.svelte', () => {
  const mockDevices = [
    { id: 'wled-1', name: 'Strip 1', ledCount: 60 },
    { id: 'wled-2', name: 'Strip 2', ledCount: 30 }
  ];

  it('renders modal header when isOpen is true', () => {
    const { getByText } = render(CreateGroupModal, {
      isOpen: true,
      allDevices: mockDevices
    });

    expect(getByText('Create Virtual Strip Group')).toBeDefined();
  });

  it('submits form when name typed and strip selected', async () => {
    const onCreate = vi.fn();
    const { getByPlaceholderText, getByText } = render(CreateGroupModal, {
      isOpen: true,
      allDevices: mockDevices,
      onCreate
    });

    const input = getByPlaceholderText('e.g. Desk Lights, TV Setup');
    await fireEvent.input(input, { target: { value: 'Desk Setup' } });

    const strip1 = getByText('Strip 1');
    await fireEvent.click(strip1);

    const submitBtn = getByText('Create Group');
    await fireEvent.click(submitBtn);

    expect(onCreate).toHaveBeenCalledWith('Desk Setup', '', ['wled-1']);
  });
});
