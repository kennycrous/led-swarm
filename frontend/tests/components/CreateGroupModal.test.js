import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import CreateGroupModal from '../../src/lib/components/CreateGroupModal.svelte';

describe('CreateGroupModal.svelte', () => {
  const mockDevices = [
    { id: 'wled-1', name: 'Strip 1', ledCount: 60 },
    { id: 'wled-2', name: 'Strip 2', ledCount: 30 }
  ];

  it('renders group creation title and device checkboxes', () => {
    const { getByText } = render(CreateGroupModal, {
      props: {
        isOpen: true,
        allDevices: mockDevices
      }
    });

    expect(getByText('Create Virtual Strip Group')).toBeDefined();
    expect(getByText('Strip 1')).toBeDefined();
    expect(getByText('Strip 2')).toBeDefined();
  });

  it('submits form with group name and selected device IDs', async () => {
    const handleSubmit = vi.fn();
    const { getByPlaceholderText, getByText } = render(CreateGroupModal, {
      props: {
        isOpen: true,
        allDevices: mockDevices,
        onCreate: handleSubmit
      }
    });

    const input = getByPlaceholderText('e.g. Living Room TV Zone');
    await fireEvent.input(input, { target: { value: 'Desk Setup' } });

    const strip1Btn = getByText('Strip 1');
    await fireEvent.click(strip1Btn);

    const submitBtn = getByText('Save Group');
    await fireEvent.click(submitBtn);

    expect(handleSubmit).toHaveBeenCalledWith('Desk Setup', '', ['wled-1']);
  });
});
